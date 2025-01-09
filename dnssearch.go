package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/miekg/dns"
)

func queryDNS(recordtype uint16, fqdn string, client *dns.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	var msg dns.Msg
	msg.SetQuestion(fqdn, recordtype)
	resp, _, err := client.Exchange(&msg, "1.1.1.1:53")
	if err != nil {
		return
	}
	
	for _, answer := range resp.Answer {
		fmt.Println(answer)
	}
}

func dnsrecords() {
	domain := ""
	fmt.Print("Please specify a Hostname or IP Address (Ex: google.com): ")
	fmt.Scan(&domain)

	if domain == "" {
		fmt.Println("Domain name cannot be empty!")
		return
	}

	fqdn := dns.Fqdn(domain)

	recordTypes := []uint16{
		dns.TypeA, dns.TypeAAAA, dns.TypeCNAME, dns.TypeMX,
		dns.TypeNS, dns.TypePTR, dns.TypeSOA, dns.TypeSRV,
		dns.TypeTXT, dns.TypeSPF, dns.TypeAAAA, dns.TypeDNSKEY,
		dns.TypeDS, dns.TypeNAPTR, dns.TypeRRSIG, dns.TypeANY,
	}

	client := dns.Client{
		Timeout: time.Second * 5,
	}

	var wg sync.WaitGroup

	for _, recordtype := range recordTypes {
		wg.Add(1)
		go queryDNS(recordtype, fqdn, &client, &wg)
	}

	wg.Wait()
}
