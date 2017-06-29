package main

/* 
	sample app, chating application interchanging text message.
*/

import(
	"fmt"
	"github.com/rainer37/Rnet/transport"
)

var app_id int = 1000

func main() {
	fmt.Println("Sample Chatting App Version 1.0 started")
	
	for {
		var cmd string
		fmt.Scanf("%s", &cmd)
		
		transport.Send_UDS("", cmd)
	}
}