package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func subdirectoryBrute() {
	hostname := ""
	wordlist := ".\\wordlists\\subdirectory.txt" 
	fmt.Print("Enter Target (google.com): ")
	fmt.Scan(&hostname)

	if hostname == "" {
		fmt.Println("Empty Hostname")
		return
	}
	userwordlist:=""

	fmt.Print("Enter wordlist (Press Enter to use default): ")
	fmt.Scan(&userwordlist)

	if userwordlist != "" {
		userwordlist = wordlist
	}

	outputFile, err := os.Create(fmt.Sprintf("%s-target.txt", hostname))
	if err != nil {
		fmt.Println("Could not create output file:", err)
		return
	}
	defer outputFile.Close()

	resultChan := make(chan string)
	var wg sync.WaitGroup

	file, err := os.Open(userwordlist)
	if err != nil {
		fmt.Println("Could not open wordlist:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subdirectory := scanner.Text()

		wg.Add(1)
		go checkSubdirectory(subdirectory, hostname, &wg, resultChan)
	}

	go func() {
		for result := range resultChan {
			fmt.Println(result)
			_, err := outputFile.WriteString(result + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}
			fmt.Println(result)
		}
	}()

	wg.Wait()
	close(resultChan)

	if err := scanner.Err(); err != nil {
		fmt.Println("Wordlist Error: ", err)
	}
}

func checkSubdirectory(subdirectory, domain string, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done()

	fullDomain := fmt.Sprintf("http://%s/%s", domain, subdirectory)
	client := http.Client{
		Timeout: 5 * time.Second, 
	}

	resp, err := client.Get(fullDomain)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		resultChan <- fullDomain
	}
}


