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
)

const TRANS_PREFIX string = "[TNS]\t"
var SOCKET_ADDR string = os.Getenv("PWD")+"/go.sock"

/*
	target_addr : remote peer address (IP/?)

	send message over UDS to DHT node and to peer.
*/
func Send_to_UDS(target_addr string, msg string) {

	if target_addr == "" {
		target_addr = SOCKET_ADDR // default
	}

	conn, err := net.Dial("unix", target_addr)

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




