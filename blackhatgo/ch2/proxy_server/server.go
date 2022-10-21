package main

// This is a simple echo server from the Black Hat Go book
// It will listen on port 20080 and echo back any data sent by connected clients

import (
	"log"
	"net"
	"io"
	"strings"
)

// echo is a handler function that simply echoes received data
func echo(conn net.Conn) {
	defer conn.Close()

	clientIP := strings.Split(conn.RemoteAddr().String(), ":")[0]

	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatalln("Unable to read/write data from ", clientIP)
	} 

	log.Printf("Client disconnected: %s", clientIP)
}

func main() {
	// Bind to TCP port 20080 on all interfaces
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}

	log.Println("Listening on 0.0.0.0:20080")

	for {
		// Wait for connection. Create net.Conn on connection established
		conn, err := listener.Accept()
		clientIP := strings.Split(conn.RemoteAddr().String(), ":")[0]
		log.Println("Received connection from: ", clientIP)
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}

		// Handle the connection.  Using goroutine for concurrency
		go echo(conn)
	}
}