package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"

	"bing-metadata-scraper/metadata"
)

var counter int = 0

func handler(i int, s *goquery.Selection) {
	url, ok := s.Find("a").Attr("href")
	if !ok {
		return
	}

	// First request retrieves an in between page with final URL in javascript code
	res, err := http.Get(url)
	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	defer res.Body.Close()

	// Regexp to extract final URL of the actual document from javascript
	urlMatch := `var u = "(.*)"`
	re := regexp.MustCompile(urlMatch)
	matches := re.FindStringSubmatch(string(buf))
	finalUrl := matches[1]

	// Make 2nd request to get the actual document
	fmt.Printf("%d: %s\n", counter, finalUrl)
	counter++
	res2, err := http.Get(finalUrl)
	if err != nil {
		return
	}

	// Read file data into buffer
	buf2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		return
	}
	defer res2.Body.Close()

	// Create zip object from file buffer
	r, err := zip.NewReader(bytes.NewReader(buf2), int64(len(buf2)))
	if err != nil {
		return
	}

	// Use defined structs and functions to unzip file and extract metadata properties
	cp, ap, err := metadata.NewProperties(r)
	if err != nil {
		return
	}

	log.Printf(
		"%25s %25s - %s %s\n",
		cp.Creator,
		cp.LastModifiedBy,
		ap.Application,
		ap.GetMajorVersion())
}

// Helper function to make request to a given URL and return a goquery document
func makeRequest(searchUrl string) (document *goquery.Document) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Request now seems to require a valid user agent or it returns no results
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0")
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panicln(err)
	}

	return doc
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Missing required argument. Usage: main.go domain ext")
	}
	domain := os.Args[1]
	filetype := os.Args[2]

	var totalUrls int
	urlCount := 0

	// Modified dork query as the instreamset option was returning 0 results
	q := fmt.Sprintf(
		"site:%s && filetype:%s",
		domain,
		filetype)
	search := fmt.Sprintf("http://www.bing.com/search?q=%s&first=%s", url.QueryEscape(q), strconv.Itoa(urlCount))

	initialDoc := makeRequest(search)

	// Regex to match the number indicating how many results there are
	urlMatch := `.*About (.+) results.*`
	re := regexp.MustCompile(urlMatch)

	// Get total number of results from query
	sTotal := "html body div#b_content div#b_tween"
	initialDoc.Find(sTotal).Each(func(i int, s *goquery.Selection) {
		// Parse the span item containing the "About X results" text
		resultsText := s.Find("span").Text()
		anyResults, _ := regexp.MatchString(urlMatch, resultsText)
		if anyResults {
			totalMatch := re.FindStringSubmatch(resultsText)[1]
			total, err := strconv.Atoi(totalMatch)
			totalUrls = total
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			fmt.Println("No results found")
			os.Exit(1)
		}

	})

	// Get total number of URLs on the first page
	sUrls := "html body div#b_content ol#b_results li.b_algo h2"
	initialDoc.Find(sUrls).Each(func(i int, s *goquery.Selection) {
		_, ok := s.Find("a").Attr("href")
		if !ok {
			return
		}
		urlCount++
	})

	// Retrieve metadata for URLs from the first page of results
	s := "html body div#b_content ol#b_results li.b_algo h2"
	initialDoc.Find(s).Each(handler)

	fmt.Printf("Total URLs seen so far: %v/%v\n", urlCount, totalUrls)

	// Loop through any additional pages of results, if they exist
	for urlCount < totalUrls {

		search = fmt.Sprintf("http://www.bing.com/search?q=%s&first=%s", url.QueryEscape(q), strconv.Itoa(urlCount))
		fmt.Printf("URL: %v\n", search)
		nextDoc := makeRequest(search)

		sUrls := "html body div#b_content ol#b_results li.b_algo h2"
		nextDoc.Find(sUrls).Each(func(i int, s *goquery.Selection) {
			_, ok := s.Find("a").Attr("href")
			if !ok {
				return
			}
			urlCount++
		})

		s := "html body div#b_content ol#b_results li.b_algo h2"
		nextDoc.Find(s).Each(handler)

		fmt.Printf("Total URLs seen on additional pages: %v/%v\n", urlCount, totalUrls)
	}

}
