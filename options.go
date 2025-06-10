package main

import "fmt"

func options() {
	option := 0
	fmt.Println("1-)Port Scan")
	fmt.Println("2-)Brute Force")
	fmt.Println("3-)Find Subdomain via Brute Force")
	fmt.Println("4-)Find Subdirectoies via Brute Force")
	fmt.Println("5-)Find Subdirectories via Scraping")
	fmt.Println("6-)Shodan Search")
	fmt.Println("7-)Dns Zone Transfer")
	fmt.Println("8-)DNS Records")
	fmt.Println("9-)AD Password Spraying")
	fmt.Println("10-)Ransomware Attack")
	fmt.Println("11-)XSS FUZZER")
	fmt.Println("12-)SQLi FUZZER")
	fmt.Print("Select your Option : ")
	fmt.Scan(&option)
	switch option {
	case 1:
		portscan()
	case 2:
		bruteforce()
	case 3:
		subdomainBrute()
	case 4:
		subdirectoryBrute()
	case 5:
		scraping()
	case 6:
		shodanoptions()
	case 7:
		zonetransfer()
	case 8:
		dnsrecords()
	case 9:
		passwordspray()
	case 10:
		ransomware()
	case 11:
		xss()
	case 12:
		sqli()
	}
}
