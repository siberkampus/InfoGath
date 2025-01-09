package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func subdomainBrute() {
	hostname := ""
	defaultWordlist := ".\\wordlists\\subdomain.txt"
	var wordlist string

	fmt.Print("Enter Target (google.com): ")
	fmt.Scan(&hostname)

	fmt.Print("Enter wordlist (for default wordlist press Enter): ")
	fmt.Scanln(&wordlist)

	if wordlist == "" {
		wordlist = defaultWordlist
	}

	if hostname == "" {
		fmt.Println("Empty Hostname")
		return
	}

	outputFileName := fmt.Sprintf("domain-%s.txt", hostname, )
	file, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	resultChan := make(chan string)
	var wg sync.WaitGroup

	wordlistFile, err := os.Open(wordlist)
	if err != nil {
		fmt.Println("Wordlist dosyası açılamadı:", err)
		return
	}
	defer wordlistFile.Close()

	scanner := bufio.NewScanner(wordlistFile)
	for scanner.Scan() {
		subdomain := scanner.Text()

		wg.Add(1)
		go checkSubdomain(subdomain, hostname, &wg, resultChan)
	}

	go func() {
		for result := range resultChan {
			fmt.Println(result)
			_, err := file.WriteString(result + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}
		}
	}()

	wg.Wait()
	close(resultChan)

	if err := scanner.Err(); err != nil {
		fmt.Println("Wordlist Error : ", err)
	}
}

func checkSubdomain(subdomain, domain string, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done()

	fullDomain := fmt.Sprintf("http://%s.%s", subdomain, domain)
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

