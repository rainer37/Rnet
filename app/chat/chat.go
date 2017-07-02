package main

/* 
	sample app, chating application interchanging text message.
*/

import(
	"fmt"
	"net"
	"github.com/rainer37/Rnet/rapi"
)

var app_id int = 1000

func echo_handler(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 1024)
	nr, err := c.Read(buf)
	if err != nil {
		return
	}

	data := buf[0:nr]
	fmt.Println("Received", string(data))
}

func receive() {
	for {
		rapi.Serve(echo_handler)
	}
}

func main() {
	fmt.Println("Sample Chatting App Version 1.0 started")
	
	rapi.Peers()

	go receive()

	for {
		var cmd, ip string
		fmt.Scanf("%s", &cmd)
		fmt.Scanf("%s", &ip)

		rapi.Send_rip(ip, cmd)
	}
}