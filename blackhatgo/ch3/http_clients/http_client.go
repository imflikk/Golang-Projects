package main

import (
	"net/http"
	"net/url"
	"strings"
	"log"
	"fmt"
	"io"
	"os"
)

func main() {
	var target string = os.Args[1]

	r1, err := http.Get(target)
	if err != nil {
		log.Fatalln("Error with GET: ", err)
	}
	// Read response body.  Not shown
	defer r1.Body.Close()
	b, err := io.ReadAll(r1.Body)

	// Print first request
	log.Println("\n\n-----GET-----\n")
	fmt.Println("Status: " + r1.Status + "\n")
	fmt.Printf(string(b))
	
	r2, err := http.Head(target)
	if err != nil {
		log.Fatalln("Error with HEAD: ", err)
	}
	// Read response body.  Not shown
	defer r2.Body.Close()
	//b, err = io.ReadAll(r1.Header)

	// Print second request
	log.Println("\n\n-----HEAD-----\n")
	fmt.Println("Status: " + r1.Status + "\n")
	for key, element := range r2.Header {
		fmt.Println(key + ": " + element[0])
	}

	form := url.Values {}
	form.Add("foo", "bar")

	r3, err := http.Post(
		target,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		log.Fatalln("Error with POST: ", err)
	}
	// Read response body.  Not shown
	defer r3.Body.Close()
	b, err = io.ReadAll(r3.Body)
	log.Println("\n\n-----POST-----\n")
	fmt.Println("Status: " + r1.Status + "\n")

	for key, element := range r3.Header {
		fmt.Println(key + ": " + element[0])
	}

	
	fmt.Printf("\n\n" + string(b) + "\n")
}