package dht

/*
	Network Scanner typically looking for any peer using R-net
	(Probably the last part to implement and it can be hard)

	Subnet_scan() 
		Scan an display all the peers using R-net if any.

*/

import (
	"fmt"
	"net"
	"strings"
	"regexp"
)

const LOCALHOST = "127.0.0.1"

func Subnet_scan() {
	fmt.Println("Looking for you...")
}

// get local ipv4 in the subnet
// TODO: get external ip 
func Local_ip_4() string {
	// list all the address from network interfaces
	faces, _ := net.InterfaceAddrs()

	for _,v := range faces {
		ip := v.String()[:strings.Index(v.String(), "/")]

		// check for ipv4 format
		if m,_ := regexp.MatchString(`\d+\.\d+\.\d+\.\d+`, ip); m && ip != LOCALHOST{
			fmt.Println(ip)
			return ip
		}
	}

	return ""
}

// get local MAC address
// TODO: what if there is more than one
func Local_MAC() string {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.Name == "en0" {
	         return inter.HardwareAddr.String()
		}
	}
	fmt.Println("No MAC found.")
	return ""
}