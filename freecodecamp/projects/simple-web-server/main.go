package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func logHTTPRequest(r *http.Request) {

	fmt.Printf("[*] %s - %s - %s - %s\n", r.RemoteAddr, r.Method, r.URL, r.Header["User-Agent"])
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	logHTTPRequest(r)

	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	homeData, err := os.ReadFile("./static/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprint(w, string(homeData))

}

func formPageHandler(w http.ResponseWriter, r *http.Request) {

	logHTTPRequest(r)

	if r.URL.Path != "/form.html" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	homeData, err := os.ReadFile("./static/form.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprint(w, string(homeData))

}

func formHandler(w http.ResponseWriter, r *http.Request) {

	logHTTPRequest(r)

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	fmt.Fprintf(w, "POST request successful")
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "\n\nName: %s\n", name)
	fmt.Fprintf(w, "\nAddress: %s\n", address)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	logHTTPRequest(r)

	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello works!")

}

func main() {
	//fileServer := http.FileServer(http.Dir("./static"))

	// Probably a better way to get logging than adding functions for each endpoint
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/form.html", formPageHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("[+] Starting web server on port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
