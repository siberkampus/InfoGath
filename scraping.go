package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func scraping() {
	url := ""
	fmt.Print("Enter Target (ex:https://apple.com): ")
	fmt.Scan(&url)
	if url == "" {
		fmt.Println("Target is empty")
		return
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var externalLinks []string
	var internalLinks []string

	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		link, exists := item.Attr("href")
		if exists {
			if strings.HasPrefix(link, "http") {
				externalLinks = append(externalLinks, link)
			} else {
				internalLinks = append(internalLinks, link)
			}
		}
	})

	outputFile, err := os.Create("subdirectory-target.txt")
	if err != nil {
		log.Fatal("Could not create output file:", err)
	}
	defer outputFile.Close()

	if len(externalLinks) > 0 {
		outputFile.WriteString("External Links:\n")
		for _, link := range externalLinks {
			outputFile.WriteString(link + "\n")
		}
	}

	if len(internalLinks) > 0 {
		outputFile.WriteString("Internal Links:\n")
		for _, link := range internalLinks {
			outputFile.WriteString(link + "\n")
		}
	}

	fmt.Println("External Links:")
	for _, link := range externalLinks {
		fmt.Println(link)
	}

	fmt.Println("Internal Links:")
	for _, link := range internalLinks {
		fmt.Println(link)
	}
}

