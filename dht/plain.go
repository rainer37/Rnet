package dht

/*
	Unstructed Plain Node Overlays

	OpStatus:
		0 : Success!
		1 : Failure!
*/

import (
	"fmt"
	"net"
	//"os"
)

const PLAIN_TOPO string = "PLAIN_OVERLAY"

type Plain_node struct {
	IP string
	Port_int int
	Port_string string
	ID string
}

func eprint(err error)  {
	fmt.Println(DHT_PREFIX, err.Error())
}

func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	defer conn.Close()

	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)

	if err != nil {
		eprint(err)
		return 
	}

	// Send a response back to person contacting us.
	conn.Write([]byte("Message received: "+string(buf)))
	// Close the connection when you're done with it.
}

func (p *Plain_node) Init() (uint8, error) {
	
	fmt.Printf("%sMy IP is [%s:%d] using [%s] as base arch\n", DHT_PREFIX, p.IP, p.Port_int, PLAIN_TOPO)
	
	p.ID = "0" // first node in the system.

	l, err := net.Listen("tcp", p.IP+":"+p.Port_string)
	if err != nil {
		eprint(err)
		return 1, err
	}

	defer l.Close()

	fmt.Println(DHT_PREFIX+"Start Listening on ["+p.IP+":"+p.Port_string+"]")

	for {
		conn, err := l.Accept()

		fmt.Println(DHT_PREFIX+"Incoming connection from "+conn.RemoteAddr().String())

		if err != nil {
			eprint(err)
			conn.Close()
			return 1, err
		}
		go handleRequest(conn)
	}

	return 0,nil
}

func (p *Plain_node) Join(ip string, port int) (uint8, error) {
	return 0, nil
} 