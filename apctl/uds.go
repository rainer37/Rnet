package apctl
/*
	UDS functions
*/

import(
	"fmt"
	"net"
	"os"
)

/*
	uds_addr : socket file name
	handler  : connection handler

	Listen for local application/DHT communication data.
*/
func Listen_UDS(uds_addr string, handler func(net.Conn)) {
	fmt.Println(AC_PREFIX+"Starting UDS server")
	
	// remove previous socket.
	os.Remove(uds_addr)

	ln, err := net.Listen("unix", uds_addr)
	global_ln = ln

	if err != nil {
		fmt.Println(AC_PREFIX+"Listen error: ", err)
		return
	}

	defer ln.Close()

	for {
		sock, err := ln.Accept()

		if err != nil {
			fmt.Println(AC_PREFIX+"Accept error: ", err)
			break
		}

		go handler(sock)
	}
}

// close the UDS listener.
func Close_apc_uds_ln() {
	if global_ln != nil {
		global_ln.Close()
	}
	fmt.Println(AC_PREFIX+"UDS Listener Closed.")
}

func Send_to_UDS(target_addr string, msg string) {

	if target_addr == "" {
		target_addr = SOCKET_ADDR // default
	}

	conn, err := net.Dial("unix", target_addr)

	if err != nil {
		fmt.Println(AC_PREFIX+"UDS Dial error", err)
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(AC_PREFIX+"Send error:", err)
		return
	}
	fmt.Println(AC_PREFIX+"Sending:", msg)

}