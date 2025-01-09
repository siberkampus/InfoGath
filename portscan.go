package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

func portscan() {
	hostname := ""
	startPort := 0
	endPort := 0
	fmt.Print("Enter target (ex:google.com): ")
	fmt.Scan(&hostname)
	fmt.Print("Enter startPort: ")
	fmt.Scan(&startPort)
	fmt.Print("Enter endPort: ")
	fmt.Scan(&endPort)
	if hostname == "" {
		fmt.Println("Hostname is blank")
		return
	}
	if startPort < 1 || startPort > 65535 {
		fmt.Println("Invalid start port")
		return
	}
	if endPort < 1 || endPort > 65535 {
		fmt.Println("Invalid end port")
		return
	}
	if startPort > endPort {
		fmt.Println("Start port cannot be higher than end port")
		return
	}


	fileName := fmt.Sprintf("port-scan+%s.txt", hostname)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex
	openPorts := []int{}
	log.Println("Start Time")
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go scanPort("tcp", hostname, port, &wg, &openPorts, &mu)
	}
	wg.Wait()

	sort.Ints(openPorts)

	for _, port := range openPorts {
		result := fmt.Sprintf("Open Port: %d\n", port)
		fmt.Print(result)
		_, err := file.WriteString(result)
		if err != nil {
			log.Fatalf("Failed to write to file: %s", err)
		}
	}
	log.Println("End Time")
}

func scanPort(protocol, hostname string, port int, wg *sync.WaitGroup, openPorts *[]int, mu *sync.Mutex) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", hostname, port)
	success := false

	for attempt := 0; attempt < 2; attempt++ { 
		conn, err := net.DialTimeout(protocol, address, time.Second*3)
		if err != nil {
			continue
		}
		conn.Close()
		success = true
		break
	}

	if success {
		mu.Lock()
		*openPorts = append(*openPorts, port)
		mu.Unlock()
	}
}

