package main

// Basic tcp-proxy server to relay traffic from one host to another.
// Doesn't currently work with encrypted traffic (SSL/TLS)

// Example: User connects to myproxy.com:8080 (where this is running)
// traffic is then forwarded to the defined host, here yahoo.com:80

import (
	"net"
	"log"
	"io"
)

func handle(src net.Conn) {
	dst, err := net.Dial("tcp", "yahoo.com:80")
	if err != nil {
		log.Fatalln("Unable to connect to our unreachable host")
	}

	defer dst.Close()

	go func() {
		// Copy our source's output to the destination
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()

	// Copy our destination's output back to our source
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	// Listen on local port 80
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	} else {
		log.Printf("Listening on port 8080\n")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}

		go handle(conn)
	}
}