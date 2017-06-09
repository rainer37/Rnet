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
	LEAVE
	REPORT
	.
	.
	.
	
*/

import (
	"fmt"
)

const DHT_PREFIX string = "[DHT]\t"

// DHT node interface
type D_node interface {
	
	/*
		Self-initialization of the new p2p systems.
	*/
	Init() (uint8, error)

	/*
		After finding the existing peer(s) in the system, send out 
		the request for joining the system.

 		Join(ip string, port int) (opStatus int, err error)
	*/
	Join(string, int) (uint8, error)

	/*
		Send the data to the peer

		Send(ip string, port int, msg []byte) (opStatus int, err error)
	*/

	// Send(string, int, []byte) (int, error)

}

// First time initialization
func Self_init(d_node D_node) error {
	
	fmt.Println(DHT_PREFIX+"General Dispatcher Initiated")
	d_node.Init()
	return nil

}