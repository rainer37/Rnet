package main

/* 
	sample app, chating application interchanging text message.
*/

import(
	"fmt"
	"net"
	"github.com/rainer37/Rnet/transport"
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
		transport.Listen_UDS("chat.sock", echo_handler)
	}
}

func main() {
	fmt.Println("Sample Chatting App Version 1.0 started")
	
	for {
		var cmd string
		fmt.Scanf("%s", &cmd)
		
		transport.Send_UDS("", cmd+":chat.sock")
	}
}