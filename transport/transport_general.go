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
)

const SOCKET_ADDR = "/tmp/go.sock"
const TRANS_PREFIX string = "[TRANS] "

type UDS struct {
	Socket_adr string // socket path.
	App_name string // binding application's name
}

func Trans_Boot() {
	fmt.Println(TRANS_PREFIX+"Transport Communication Unit Booting...")
	Listen_UDS()
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
	Listen for local application communication data.
*/
func Listen_UDS() {
	fmt.Println(TRANS_PREFIX+"Starting UDS server")
	
	ln, err := net.Listen("unix", SOCKET_ADDR)

	if err != nil {
		fmt.Println(TRANS_PREFIX+"Listen error: ", err)
		return
	}

	defer func (ln net.Listener) {
		fmt.Println(TRANS_PREFIX+"Closing Socket")
		ln.Close()
	}(ln)

	for {
		sock, err := ln.Accept()
		
		if err != nil {
			fmt.Println(TRANS_PREFIX+"Accept error: ", err)
			break
		}

		go handle_UDS_request(sock)
	}
}

func handle_UDS_request(c net.Conn) {
		defer c.Close()
	//for {
		buf := make([]byte, 1024)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		fmt.Println("Server got:", string(data))
		// _, err = c.Write(data)
		// if err != nil {
		// 	fmt.Println("Writing client error: ", err)
		// }
	//}
}
