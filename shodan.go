package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Count struct {
	Total int `json:"total"`
}

func shodanIpScan() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not set in environment")
	}
	target := ""
	fmt.Print("Enter target IP address :")
	fmt.Scan(&target)
	requesturl := fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", target, apiKey)

	res, err := http.Get(requesturl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Shodan API request failed with status: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response ShodanResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("IP: %s\n", response.IPStr)
	fmt.Printf("Ports: %d\n", response.Ports[:])
	fmt.Printf("CountryCode: %v\n", response.CountryCode)
}
func shodanoptions() {
	option := 0
	fmt.Println("1-)Shodan IP Scan")
	fmt.Println("2-)Shodan Count")
	fmt.Println("3-)Get Shodan Filters")
	fmt.Println("4-)Shodan Search Filter")
	fmt.Print("Select your Option : ")
	fmt.Scan(&option)
	switch option {
	case 1:
		shodanIpScan()
	case 2:
		shodanCount()
	case 3:
		shodanFilters()
	case 4:
		shodanResult()
	}
}
func shodanCount() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not set in environment")
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your query: ")
	reader.ReadString('\n')
	query, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading input:", err)
	}
	query = strings.TrimSpace(query)

	encodedquery := url.QueryEscape(query)
	requesturl := fmt.Sprintf("https://api.shodan.io/shodan/host/count?key=%s&query=%s", apiKey, encodedquery)
	res, err := http.Get(requesturl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Shodan API request failed with status: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response Count
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total %d target found", response.Total)
}

func shodanFilters() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not set in environment")
	}

	requesturl := fmt.Sprintf("https://api.shodan.io/shodan/host/search/filters?key=%s", apiKey)
	res, err := http.Get(requesturl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Shodan API request failed with status: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}

func shodanResult() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not set in environment")
	}
	fmt.Scanln()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your query: ")
	query, _ := reader.ReadString('\n')
	fmt.Println("Query:", query)
	query = url.QueryEscape(query)
	requesturl := fmt.Sprintf("https://api.shodan.io/shodan/host/search?key=%s&query=%s", apiKey, query)
	fmt.Println(requesturl)
	res, err := http.Get(requesturl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Shodan API request failed with status: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response SearchResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := os.Create("result.txt")
	if err != nil {
		log.Fatal("Error creating result.txt:", err)
	}
	defer outputFile.Close()

	for _, service := range response.Matches {
		result := fmt.Sprintf("\n\nIP: %s\n", service.HostInfo.IPstr)
		result += fmt.Sprintf("City: %s\n", *service.Location.City)
		result += fmt.Sprintf("Country: %s\n", *service.Location.CountryName)
		//result += fmt.Sprintf("Data: %s\n", service.Data)
		result += fmt.Sprintf("Port: %d\n", service.Port)
		result += fmt.Sprintf("Domains: %v\n", service.Domains)

		if service.Vulns != nil {
			result += "Vulnerabilities:\n"
			count := 0
			for vulnID, vuln := range service.Vulns {
				result += fmt.Sprintf("\tVulnerability ID: %s\n", vulnID)
				result += fmt.Sprintf("\tCVSS Score: %v\n", vuln.CVSS)
				result += fmt.Sprintf("\tVerified: %v\n", vuln.Verified)
				count++
				if count == 10 {
					break
				}
			}
		} else {
			result += "\tNo vulnerabilities found.\n"
		}
		_, err := outputFile.WriteString(result)
		if err != nil {
			log.Fatal("Error writing to file:", err)
		}
		fmt.Println(result)
	}
}
