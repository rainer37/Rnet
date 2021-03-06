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
*/

import (
	"os"
	"fmt"
	"net"
	"errors"
)

// basic node structure.
type Node struct {
	IP string
	State uint
	Port_string string
	ID string
	NList interface{}
}

type Peer struct {
	peer_id string // peer's unique id
	addr string
}

const	(
	PLAIN_TOPO string = "plain"
	CHORD_TOPO string = "chord"
	DHT_PREFIX string = "[DHT]\t"
)

var SOCKET_ADDR string = os.Getenv("PWD")+"/go.sock"
var peers map[string]Peer // potential peers maybe

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

	Send(net.Conn, string) (uint8, error)

}

// First time initialization
func Self_init(ip string, port string, topo string) error {

	fmt.Println(DHT_PREFIX+"Main DHT Dispatcher Initiated")

	// dispatch the actual topo node
	var d_node D_node = dispatch(topo)

	if d_node == nil {
		fmt.Println(DHT_PREFIX+"Unknown Topology Provided, System Exiting...")
		return errors.New("Unknown Topology...")	
	}

	// start initialize the network.
	_, err := d_node.Init(ip, port)
	if err != nil {
		return err
	}

	return nil
}

// auto join or liveness check for the returning user.
func Auto_join(states string) error {
	return nil
}

// Join wrapper
func Want_to_join(ip string, port string, states string, my_port string) error {

	// TODO: finish loadfromstates and pop topo from it to dispatch.
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

// communication between DHT nodes.
// ??? why do i even need this level of abstraction
func Send_to_ext_DHT(target_addr string, msg string) {
	node := Get()
	fmt.Println(node.IP)

	conn, err := net.Dial("tcp", target_addr)
	defer conn.Close()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	conn.Write([]byte("U "+msg))
}

func load_from_states_string(states string) {

}

// dispatcher for d_node.
func dispatch(topo string) (d_node D_node) {
	switch topo {
	case PLAIN_TOPO:
		d_node = new(Plain_node)
	case CHORD_TOPO:
		//s_node := new(Chord_node)
	default:
		fmt.Println(DHT_PREFIX+"Cannot dispatch!")
		d_node = nil
	}
	return
}

/*
	forward raw internet data to application UDS.
*/
func forward_to_app(msg string) {
	conn, err := net.Dial("unix", SOCKET_ADDR)

	if err != nil {
		fmt.Println(DHT_PREFIX+"UDS Dial error", err)
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(DHT_PREFIX+"Send error:", err)
		return
	}
}
