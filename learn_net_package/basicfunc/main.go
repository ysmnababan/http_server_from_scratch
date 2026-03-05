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

	_, addrs, _ := net.LookupSRV("xmppp-server", "tcp", "google.com")
	fmt.Println("SRV records for xmppp-server.google.com:", addrs)

	ifs, _ := net.Interfaces()
	for _, iface := range ifs {
		fmt.Printf("Interface: %s, Flags: %v\n", iface.Name, iface.Flags)

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			fmt.Printf("  Address: %v\n", addr)
		}
	}

	// concurrent TCP server
	ln, _ := net.Listen("tcp", ":8080")
	for {
		conn, _ := ln.Accept()

		go func(c net.Conn) {
			defer c.Close()
			b := make([]byte, 1024)
			for {
				n, err := c.Read(b)
				if err != nil {
					fmt.Println("client disconnected")
					return
				}

				c.Write(b[:n])
			}
		}(conn)
	}
}
