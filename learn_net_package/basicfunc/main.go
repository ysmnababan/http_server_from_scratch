package main

import (
	"fmt"
	"net"
)

func main() {
	// IP tools
	ip := net.ParseIP("192.168.1.1")
	fmt.Println("Parsed IP:", ip)

	ip, ipnet, _ := net.ParseCIDR("192.168.1.2/16")
	fmt.Printf("Parsed CIDR: IP=%v, Network=%v\n", ip, ipnet)

	if ipnet.Contains(net.ParseIP("192.169.12.50")) {
		fmt.Println("IP is in the network")
	}

	// DNS
	ips, _ := net.LookupIP("www.google.com")
	fmt.Println("Google IPs:", ips)

	cname, _ := net.LookupCNAME("www.google.com")
	fmt.Println("CNAME for www.google.com:", cname)

	mx, _ := net.LookupMX("gmail.com")
	for _, record := range mx {
		fmt.Println("MX records for gmail.com:", record)
	}
}
