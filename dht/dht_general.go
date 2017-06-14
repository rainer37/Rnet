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
	"os"
	"fmt"
	"net"
)

// basic node structure.
type Node struct {
	IP string
	State uint
	Port_string string
	ID string
}

const	(
	PLAIN_TOPO string = "plain"
	CHORD_TOPO string = "chord"
	DHT_PREFIX string = "[DHT]\t"
)

// DHT node interface
type D_node interface {
	
	/*
		Self-initialization of the new p2p systems.
		Init(ip string, port string, extra ...string) (opStatus uint8, err error)
	*/
	Init(string, string, ...string) (uint8, error)

	/*
		After finding the existing peer(s) in the system, send out 
		the request for joining the system.

 		Join(ip string, port int) (opStatus uint8, err error)
	*/
	Join(string, string) (uint8, error)

	/*
		Send the data to the peer

		Send(conn net.Conn, msg string) (opStatus int, err error)
	*/

	Send(net.Conn, string) (int, error)

}

// First time initialization
func Self_init(ip string, port string, topo string) error {

	fmt.Println(DHT_PREFIX+"Main DHT Dispatcher Initiated")
	load_from_states_string("")

	switch topo {
	case PLAIN_TOPO:
		d_node := new(Plain_node)
		_, err := d_node.Init(ip, port)
		if err != nil {
			return err
		}
	default:
		fmt.Println(DHT_PREFIX+"Unknown Topology Provided, System Exiting...")
		os.Exit(-1)
	}
	
	return nil
}

// auto join or liveness check for the returning user.
func Auto_join(states string) error {
	return nil
}

// Join wrapper
func Want_to_join(ip string, port string, states string, my_port string) error {

	p := new(Plain_node)

	p.IP = Local_ip_4()
	p.Port_string = my_port
	p.State = 0
	p.NList = make(map[string]string)

	_, err := p.Join(ip, port)

	if err != nil {
		os.Exit(-1)
	}

	return nil
}

func load_from_states_string(states string) {

}