package main

/* 
	sample app, chating application interchanging text message.
*/

import(
	"fmt"
	"net"
)

func main() {
	fmt.Println("Sample Chatting App Version 1.0 started")
	
	laddr,_ := net.ResolveTCPAddr("tcp", "192.168.31.101:1339")
	raddr,_ := net.ResolveTCPAddr("tcp", "192.168.31.101:1338")


	for {
		var cmd string
		fmt.Scanf("%s", &cmd)
		
		conn,err := net.DialTCP("tcp", laddr, raddr)

		if err != nil {
			fmt.Println(err.Error())
			break
		}

		fmt.Println("CMD received:", cmd)
		conn.Write([]byte(cmd))
		conn.Close()
	}
}