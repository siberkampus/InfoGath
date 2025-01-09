package main

type SearchResult struct {
	Matches []Service `json:"matches"`
	Total   int       `json:"total"`
}

type ExploitResult struct {
	Matches []Exploit `json:"matches"`
	Total   int       `json:"total"`
}
type ShodanResponse struct {
	IPStr       string `json:"ip_str"`
	Ports       []int  `json:"ports"`
	CountryCode string `json:"country_code"`
}

type Service struct {
	HostInfo
	Location Location `json:"location"`

	// Contains the banner information for the service
	Data string `json:"data"`

	// The IP address of the host as an integer
	IP *int `json:"ip,omitempty"`

	// The IPv6 address of the host as a string. If this is present then the "IP" and "IPstr" fields wont be.
	IPv6 *string `json:"ipv6,omitempty"`

	// The port number that the service is operating on
	Port int `json:"port"`

	// The timestamp for when the banner was fetched from the device in the UTC timezone.
	// Example: "2014-01-15T05:49:56.283713"
	Timestamp string `json:"timestamp"`

	// Numeric hash of the data property
	Hash int `json:"hash"`

	// An array of strings containing the top-level domains for the hostnames of the device. This is a utility property
	// in case you want to filter by TLD instead of subdomain. It is smart enough to handle global TLDs with several
	// dots in the domain (ex. "co.uk")
	Domains []string `json:"domains"`

	// The network link type. Possible values are: "Ethernet or modem", "generic tunnel or VPN", "DSL", "IPIP or SIT",
	// "SLIP", "IPSec or GRE", "VLAN", "jumbo Ethernet", "Google", "GIF", "PPTP", "loopback", "AX.25 radio modem".
	Link *string `json:"link,omitempty"`

	// Contains experimental and supplemental data for the service. This can include the SSL certificate, robots.txt
	// and other raw information that hasn't yet been formalized into the Banner Specification.
	Opts map[string]interface{} `json:"opts"`

	// Uptime of the IP (in minutes)
	Uptime *int `json:"uptime,omitempty"`

	// Either "udp" or "tcp" to indicate which IP transport protocol was used to fetch the information
	Transport string `json:"transport"`

	// Name of the software running the service
	// In rare occasions can be number. Use ProductString() to get string value of Product
	Product interface{} `json:"product,omitempty"`

	// Version of the software
	// In rare occasions can be number. Use VersionString() to get string value of Version
	Version interface{} `json:"version,omitempty"`

	// Common platform enumeration
	// Can be slice of strings or single string. Use CpeList() to get a slice of strings value
	CPE interface{} `json:"cpe,omitempty"`

	// The title of the website as extracted from the HTML source
	Title *string `json:"title,omitempty"`

	// The type of device (webcam, router, etc.)
	DeviceType *string `json:"devicetype,omitempty"`

	// Miscellaneous information that was extracted about the product
	Info *string `json:"info,omitempty"`

	// The _shodan property contains information about how the banner was generated. It doesn’t store any
	// data about the port/service itself.
	// Availability: All banners
	Shodan Shodan `json:"_shodan"`

	// The vulns property contains information about vulnerabilities that may exist in the service represented
	// by the banner. In general, the Shodan crawlers don’t perform vulnerability testing as a result the
	// vulnerabilities stored in vulns are inferred from the banner and haven’t been verified.
	// Availability: Banners where the software/version has been identified and there exist known CVEs for it
	Vulns map[string]Vulnerability `json:"vulns,omitempty"`

	// Availability: Services that require SSL (ex. HTTPS) or support upgrading a connection to SSL/TLS
	// (ex. POP3 with STARTTLS)
}

type Exploit struct {
	// Unique ID for the exploit/vulnerability
	Id interface{} `json:"_id"`

	// The author of the exploit/vulnerability.
	Author interface{} `json:"author"`

	// The Bugtraq ID for the exploit.
	BID interface{} `json:"bid"`

	// The actual code of the exploit.
	Code string `json:"code"`

	// The Common Vulnerability and Exposures ID for the exploit.
	CVE []string `json:"cve"`

	// When the exploit was released.
	Date string `json:"date"`

	// The description of the exploit, how it works and where it applies.
	Description string `json:"description"`

	// The name of the data source. Possible values are: CVE, ExploitDB, Metasploit
	Source string `json:"source"`

	// The Microsoft Security Bulletin ID for the exploit.
	MSB interface{} `json:"msb"`

	// The Open Source Vulnerability Database ID for the exploit.
	OSVDB interface{} `json:"osvdb"`

	// The operating system that the exploit targets. Possible values are: aix, cgi, freebsd, hardware, Java, jsp,
	// lin_x86, Linux, multiple, novell, osx, PHP, true64, Unix, Windows
	// sometimes string, sometimes array of strings
	Platform interface{} `json:"platform"`

	// The port number for the affected service if the exploit is remote.
	Port int `json:"port"`

	// The title or short description for the exploit if available.
	Title string `json:"title"`

	// The type of exploit
	Type string `json:"type"`

	Alias      interface{} `json:"alias"`
	Rank       interface{} `json:"rank"`
	Arch       interface{} `json:"arch"`
	Privileged bool        `json:"privileged"`
	Version    interface{} `json:"version"`
}

type Host struct {
	Ports      []int      `json:"ports"`
	Vulns      []string   `json:"vulns"`
	LastUpdate string     `json:"last_update"`
	Services   []*Service `json:"data"`
	Location
	HostInfo
}

// common fields for "/host/{ip}" and "/host/search" endpoint results
type HostInfo struct {
	// The IP address of the host as a string
	IPstr string `json:"ip_str"`

	// The autonomous system number (ex. "AS4837")
	ASN *string `json:"asn,omitempty"`

	// The operating system that powers the device
	OS *string `json:"os"`

	// The name of the organization that is assigned the IP space for this device
	Org *string `json:"org"`

	// The ISP that is providing the organization with the IP space for this device.
	// Consider this the "parent" of the organization in terms of IP ownership
	ISP *string `json:"isp"`

	// An array of strings containing all of the hostnames that have been assigned to the IP address for this device
	Hostnames []string `json:"hostnames"`

	// List of tags that describe the characteristics of the device
	Tags []string `json:"tags,omitempty"`

	// Raw HTML of response
	HTML string `json:"html,omitempty"`
}

type Location struct {
	// The latitude for the geolocation of the device
	Latitude *float32 `json:"latitude"`

	// The longitude for the geolocation of the device
	Longitude *float32 `json:"longitude"`

	//  The name of the city where the device is located
	City *string `json:"city"`

	// The 2-letter country code for the device location
	CountryCode *string `json:"country_code"`

	// The 3-letter country code for the device location
	CountryCode3 *string `json:"country_code3"`

	// The name of the country where the device is located
	CountryName *string `json:"country_name"`

	// The area code for the device's location. Only available for the US
	AreaCode *int `json:"area_code"`

	// The name of the region where the device is located
	RegionCode *string `json:"region_code"`

	// The designated market area code for the area where the device is located. Only available for the US
	DmaCode *int `json:"dma_code"`

	// The postal code for the device's location
	PostalCode *string `json:"postal_code"`
}
type Shodan struct {
	// Unique ID that identifies the Shodan crawler
	Crawler string `json:"crawler"`

	// Unique ID for this banner
	Id *string `json:"id"`

	// Name of the Shodan module used by the crawler to generate this banner
	Module string `json:"module,omitempty"`

	// [NOT DOCUMENTED]
	Ptr bool `json:"ptr"`

	// Configuration options used during the data collection
	Options CrawlerOptions `json:"options"`
}

type CrawlerOptions struct {
	// Hostname to use when sending web requests
	Hostname string `json:"hostname,omitempty"`

	// ID of the banner that triggered the scan for this port
	Referrer string `json:"referrer,omitempty"`

	// ID of the scan
	Scan string `json:"scan,omitempty"`
}
type Vulnerability struct {
	// Common Vulnerability Scoring System value
	CVSS interface{} `json:"cvss"`

	// List of URLs that are related to the vulnerability
	References []string `json:"references"`

	// A description of the vulnerability
	Summary string `json:"summary"`

	// Whether or not the vulnerability has been verified by the Shodan crawlers
	Verified bool `json:"verified"`
}
  