package main

import (
	"fmt"
	"log"
	"os"

	"metasploit-minimal/rpc"
)

func main() {
	host := os.Getenv("MSFHOST")
	pass := os.Getenv("MSFPASS")
	user := "msf"
	if host == "" || pass == "" {
		log.Fatalln("Missing required environment variable MSFHOST or MSFPASS")
	}
	msf, err := rpc.New(host, user, pass)
	if err != nil {
		log.Panicln(err)
	}
	defer msf.Logout()
	sessions, err := msf.SessionList()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("Sessions:")
	for _, session := range sessions {
		fmt.Printf("%5d | %s %s/%s | %s (%v -> %v)\n", session.ID, session.Type, session.Platform, session.Arch, session.Info, session.TunnelLocal, session.TunnelPeer)
	}
}
