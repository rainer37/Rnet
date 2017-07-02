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

// close the global UDS listener and ipc dispatcher SOCKET_ADDR.
func Close_apc_uds_ln() {
	if global_ln != nil {
		global_ln.Close()
	}
	fmt.Println(AC_PREFIX+"UDS Listener Closed.")
}

func ini_uds_con(target_addr string, msg string) net.Conn {
	if target_addr == "" {
		target_addr = SOCKET_ADDR // default
	}

	conn, err := net.Dial("unix", target_addr)

	if err != nil {
		fmt.Println(AC_PREFIX+"1UDS Dial error on "+target_addr, err)
		return nil
	}

	return conn
}

/*
	send msg to uds server, either SOCKET_ADDR or specified.
	connection closes each msg sent.
*/
func Send_to_UDS(target_addr string, msg string) {

	conn := ini_uds_con(target_addr, msg)

	defer conn.Close()

	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(AC_PREFIX+"Send error:", err)
		return
	}
	fmt.Println(AC_PREFIX+"Sending:", msg)

}

func Send_and_receive(target_addr string, msg string) string {
	
	conn := ini_uds_con(target_addr, msg)

	defer conn.Close()

	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(AC_PREFIX+"Send error:", err)
		return ""
	}
	fmt.Println(AC_PREFIX+"Sending:", msg)

	buf := make([]byte, 1024)
	nr, err := conn.Read(buf)
	if err != nil {
		return ""
	}

	data := buf[0:nr]
	fmt.Println(AC_PREFIX+"Received [", string(data), "]")

	return string(data)
}