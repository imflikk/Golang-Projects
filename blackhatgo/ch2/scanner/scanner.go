package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"sort"
)

func printUsage() {
	fmt.Println("Usage Examples:\n\t./scanner <Target IP> <Port or port range>")
	fmt.Println("\t./scanner 192.168.1.1 445\n\t./scanner 192.168.1.1 1-1024")
}

func scanPorts(protocol, hostname string, ports chan int, results chan string) {
	
	// Read each port from the ports channel and try to connect to it
	// Depending on response, send a string back to the results channel
	for port := range ports {
		address := hostname + ":" + strconv.Itoa(port)
		conn, err := net.DialTimeout(protocol, address, 3*time.Second)

		if err != nil {
			// port is closed or filtered.
			sErr := err.Error()
			if strings.HasSuffix(sErr, "i/o timeout") {
				results <- strconv.Itoa(port) + "f"
				continue
			}
			results <- "closed"
			continue
		}

		conn.Close()
		results <- strconv.Itoa(port)
	}
	
	
}

func main() {

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Get arguments from the command line
	args := os.Args[1:]

	// Initialize the two channels to use for scanning
	ports := make(chan int, 100)
	results := make(chan string)

	// Initialize various variables used
	lowerPort, upperPort, port := 0, 0, 0
	openPorts := []int{}
	filteredPorts := []int{}
	target := args[0]

	start := time.Now()

	lineSeparator := "---------------------"

	// Parse command line port arguments
	if strings.Contains(args[1], "-") {
		// If port range provided, parse out starting and ending port
		portRange := strings.Split(args[1], "-")
		lowerPort, _ = strconv.Atoi(portRange[0])
		upperPort, _ = strconv.Atoi(portRange[1])
	} else {
		port, _ = strconv.Atoi(args[1])
	}

	if !strings.Contains(args[1], "-") {
		// If specific port provided
		fmt.Printf("[*] Scanning port %d on %s...\n"+lineSeparator+"\n", port, target)

		for i := 0; i < cap(ports); i++ {
			go scanPorts("tcp", target, ports, results)
		}

		go func() {
			ports <- port
		}()

	} else {
		// If port range provided
		fmt.Printf("[*] Scanning ports %d-%d on %s...\n"+lineSeparator+"\n", lowerPort, upperPort, target)

		// 
		for i := 0; i < cap(ports); i++ {
			go scanPorts("tcp", target, ports, results)
		}


		go func() {
			for currentPort := lowerPort; currentPort < upperPort; currentPort++ {
				ports <- currentPort
			}
		}()

	}

	// Read from results channel to get open/filtered ports
	for currentPort := lowerPort; currentPort < upperPort; currentPort++ {
		port := <- results
		if port != "closed" {
			if strings.HasSuffix(port, "f") {
				newPort, _ := strconv.Atoi(strings.TrimSuffix(port, "f"))
				filteredPorts = append(filteredPorts, newPort)
			} else {
				newPort, _ := strconv.Atoi(port)
				openPorts = append(openPorts, newPort)
			}
			
		}
	}

	// Close channels
	close(ports)
	close(results)


	// Sort and print out all open ports when finished
	fmt.Printf("[+] Ports open on %s:\n", target)

	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf("%d open\n", port)
	}

	//// Filtered ports isn't working reliably at the moment without a long timeout

	// sort.Ints(filteredPorts)
	// fmt.Printf(lineSeparator+"\n[*] Ports filtered on %s:\n", target)
	// for i := 0; i < len(filteredPorts); i++ {
	// 	fmt.Printf("%d\n", filteredPorts[i])
	// }

	duration := time.Since(start)
	fmt.Printf("\nTime elapsed: %.02f seconds\n", duration.Seconds())

}