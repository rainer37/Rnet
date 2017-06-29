package transport

/*
	Local communication between DHT routing unit and applications.
	Using fixed dynamic Unix Domain Socket(UDS) for all communication.
	
	General Interfaces:

	Listen_UDS()
	Send_UDS(target_addr string)

*/

import (
	"fmt"
	"net"
	"os"
	"github.com/rainer37/Rnet/dht"
)

const TRANS_PREFIX string = "[TNS]\t"

type UDS struct {
	Socket_adr string // socket path.
	App_name string // binding application's name
}

var SOCKET_ADDR string = os.Getenv("PWD")+"/go.sock"
var lnn net.Listener // pointer to listener

func Trans_Boot(id string) {
	fmt.Println(TRANS_PREFIX+"Transport Communication Unit Booting..."+id)
	Listen_UDS(SOCKET_ADDR, DHT_UDS_handler)
} 

// close the UDS listener.
func Close_ln() {
	if lnn != nil {
		lnn.Close()
	}
	fmt.Println(TRANS_PREFIX+"UDS Listener Closed.")
}

/*
	target_addr : remote peer address (IP/?)

	send message over UDS to DHT node and to peer.
*/
func Send_UDS(target_addr string, msg string) {

	conn, err := net.Dial("unix", SOCKET_ADDR)

	if err != nil {
		fmt.Println(TRANS_PREFIX+"UDS Dial error", err)
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(TRANS_PREFIX+"Send error:", err)
		return
	}
	fmt.Println(TRANS_PREFIX+"Sending:", msg)

}

/*
	Listen for local application/DHT communication data.
*/
func Listen_UDS(uds_addr string, handler func(net.Conn)) {
	fmt.Println(TRANS_PREFIX+"Starting UDS server")
	
	// remove previous socket.
	os.Remove(uds_addr)

	ln, err := net.Listen("unix", uds_addr)
	lnn = ln

	if err != nil {
		fmt.Println(TRANS_PREFIX+"Listen error: ", err)
		return
	}

	defer ln.Close()

	for {
		sock, err := ln.Accept()

		if err != nil {
			fmt.Println(TRANS_PREFIX+"Accept error: ", err)
			break
		}

		go handler(sock)
	}
}

/*
	handle message from applications 
	(and forward to DHT network for inter-DHT communication)

	TODO: application data encoding unit.
*/
func DHT_UDS_handler(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 1024)
	nr, err := c.Read(buf)
	if err != nil {
		return
	}

	data := buf[0:nr]
	sock := string(data[len(data)-8:])

	fmt.Println(TRANS_PREFIX+"Server got:", string(data), sock)

	// send data over internet DHT nodes.
	dht.Send_DHT("192.168.31.205:1338",string(data))
}
