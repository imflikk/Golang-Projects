package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func printUsage() {
	fmt.Println("Usage Examples:\n\t./scanner <Target IP or IP range> <Port or port range>")
	fmt.Println("\t./scanner 192.168.1.1 445\n\t./scanner 192.168.1.1 1-1024\n\t./scanner 192.168.1.1-192.168.1.5 80")
}

func scanPorts(protocol, hostname string, ports chan int, results chan string) {
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

// ipToInt converts a net.IP to a uint32
func ipToInt(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// intToIP converts a uint32 to a net.IP
func intToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte((n>>16)&0xFF), byte((n>>8)&0xFF), byte(n&0xFF))
}

// parseIPRange parses an IP range in the form 192.168.1.10-192.168.1.20
func parseIPRange(ipRange string) ([]string, error) {
	ips := []string{}
	if !strings.Contains(ipRange, "-") {
		// Single IP
		return []string{ipRange}, nil
	}
	parts := strings.Split(ipRange, "-")
	startIP := net.ParseIP(parts[0]).To4()
	endIP := net.ParseIP(parts[1]).To4()
	if startIP == nil || endIP == nil {
		return nil, fmt.Errorf("invalid IP range")
	}
	start := ipToInt(startIP)
	end := ipToInt(endIP)
	if start > end {
		return nil, fmt.Errorf("start IP must be less than or equal to end IP")
	}
	for i := start; i <= end; i++ {
		ips = append(ips, intToIP(i).String())
	}
	return ips, nil
}

func processTarget(target string, portArg string, wg *sync.WaitGroup, openIPPorts *[]string, mu *sync.Mutex) {
	defer wg.Done()
	ports := make(chan int, 100)
	results := make(chan string, 100)
	lowerPort, upperPort, port := 0, 0, 0
	openPorts := []int{}
	filteredPorts := []int{}
	//lineSeparator := "---------------------"

	if strings.Contains(portArg, "-") {
		// If port range provided, parse out starting and ending port
		portRange := strings.Split(portArg, "-")
		lowerPort, _ = strconv.Atoi(portRange[0])
		upperPort, _ = strconv.Atoi(portRange[1])
	} else {
		port, _ = strconv.Atoi(portArg)
	}

	fmt.Printf("\nScanning target: %s\n", target)

	if !strings.Contains(portArg, "-") {
		// If specific port provided
		//fmt.Printf("[*] Scanning port %d on %s...\n"+lineSeparator+"\n", port, target)

		for i := 0; i < cap(ports); i++ {
			go scanPorts("tcp", target, ports, results)
		}

		go func() {
			ports <- port
			close(ports)
		}()

		// Only one port to read
		portString := <-results
		if portString != "closed" {
			if strings.HasSuffix(portString, "f") {
				newPort, _ := strconv.Atoi(strings.TrimSuffix(portString, "f"))
				filteredPorts = append(filteredPorts, newPort)
			} else {
				newPort, _ := strconv.Atoi(portString)
				openPorts = append(openPorts, newPort)
				// Add to shared openIPPorts
				mu.Lock()
				*openIPPorts = append(*openIPPorts, fmt.Sprintf("%s:%d", target, newPort))
				mu.Unlock()
			}
		}

	} else {
		// If port range provided
		//fmt.Printf("[*] Scanning ports %d-%d on %s...\n"+lineSeparator+"\n", lowerPort, upperPort, target)

		for i := 0; i < cap(ports); i++ {
			go scanPorts("tcp", target, ports, results)
		}

		go func() {
			for currentPort := lowerPort; currentPort < upperPort; currentPort++ {
				ports <- currentPort
			}
			close(ports)
		}()

		for currentPort := lowerPort; currentPort < upperPort; currentPort++ {
			portString := <-results
			if portString != "closed" {
				if strings.HasSuffix(portString, "f") {
					newPort, _ := strconv.Atoi(strings.TrimSuffix(portString, "f"))
					filteredPorts = append(filteredPorts, newPort)
				} else {
					newPort, _ := strconv.Atoi(portString)
					openPorts = append(openPorts, newPort)
					// Add to shared openIPPorts
					mu.Lock()
					*openIPPorts = append(*openIPPorts, fmt.Sprintf("%s:%d", target, newPort))
					mu.Unlock()
				}
			}
		}
	}

	// Close results channel after reading all results
	close(results)

	if len(openPorts) > 0 {
		// Sort and print out all open ports when finished
		//fmt.Printf("[+] Ports open on %s:\n", target)
		sort.Ints(openPorts)
		for _, port := range openPorts {
			fmt.Printf("%d open\n", port)
		}
	}

	// Uncomment if you want to print filtered ports
	// sort.Ints(filteredPorts)
	// fmt.Printf(lineSeparator+"\n[*] Ports filtered on %s:\n", target)
	// for i := 0; i < len(filteredPorts); i++ {
	// 	fmt.Printf("%d\n", filteredPorts[i])
	// }
}

func main() {

	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	args := os.Args[1:]

	targets, err := parseIPRange(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing IP range: %v\n", err)
		os.Exit(1)
	}

	start := time.Now()
	var wg sync.WaitGroup

	var openIPPorts []string
	var mu sync.Mutex

	for _, target := range targets {
		wg.Add(1)
		go processTarget(target, args[1], &wg, &openIPPorts, &mu)
	}

	wg.Wait()

	duration := time.Since(start)
	fmt.Println("\nAll open IP:PORT combinations found:")
	sort.Strings(openIPPorts)
	for _, ipport := range openIPPorts {
		fmt.Println(ipport)
	}

	fmt.Printf("\nTime elapsed: %.02f seconds\n", duration.Seconds())
}
