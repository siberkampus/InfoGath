package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func zonetransfer() {
	domain := ""
	fmt.Print("Please specify a Hostname or IP Address (Ex: google.com) : ")
	fmt.Scan(&domain)
	if domain == "" {
		fmt.Println("Domain address is empty")
		return
	}
	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}

	client := dns.Client{
		Timeout: time.Second * 5, 
	}

	var msg dns.Msg
	fqdn := dns.Fqdn(domain) 
	msg.SetQuestion(fqdn, dns.TypeNS)

	resp, _, err := client.Exchange(&msg, "1.1.1.1:53")
	if err != nil {
		fmt.Println("Error during DNS query:", err)
		return
	}

	nameserver := ""
	for _, ans := range resp.Answer {
		if ns, ok := ans.(*dns.NS); ok {
			nameserver = ns.Ns
			nameserver = fmt.Sprintf("%s:53", strings.TrimSuffix(nameserver, ".")) 
			break
		}
	}

	if nameserver == "" {
		fmt.Println("Name server not found")
		return
	}

	transfer := new(dns.Transfer)
	msg.SetAxfr(domain) 
	channel, err := transfer.In(&msg, nameserver)
	if err != nil {
		fmt.Println("Zone transfer failed:", err)
		return
	}

	for zone := range channel {
		if zone.Error != nil {
			fmt.Println(zone.Error)
			return
		}
		for _, rr := range zone.RR {
			fmt.Println(rr)
		}
	}
}