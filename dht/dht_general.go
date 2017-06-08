package dht

/*
	R-net can apply different DHTs as required to optimize the overall performance
	This file defines the general interfaces for all DHT and the dispatcher. 
	
	INIT
	JOIN
	REPLICATE
	CHECK_PEER
	NOTIFY
	SEND
	.
	.
	.
	
*/

import (
	"fmt"
)

const DHT_PREFIX string = "[DHT]\t"

type Node struct {
	IP string
	Port_int int
	Port_string string
	ID string
}

// First time initialization
func Self_init(ip string, port int, topo string) *Node {
	fmt.Printf("%sMy IP is [%s:%d] using [%s] as base arch\n", DHT_PREFIX, ip, port, topo)
	new_node := &Node{ip, port, "", ""}

	return new_node
}