package main

/*
	User Interfaces for both client & server side ops.
*/

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"github.com/rainer37/Rnet/dht"
	"github.com/rainer37/Rnet/transport"
)

func main() {

	fmt.Println("Initiating R-NET v1.0.0\t[", time.Now() ,"]")
	fmt.Println("Local Arch & OS:\t[", runtime.GOARCH, ":", runtime.GOOS,"]")
	
	/* flags multiplexer
		-i : self initialization with [ip] [port] [topo]
		-j : join on an existing node if known ip with [target_ip] [target_port]
		TODO: MORE CASES -s -l -a
	*/
	flag, ip, port := os.Args[1],os.Args[2],os.Args[3]

	switch flag {
	case "-i":
		go dht.Self_init(ip, port, os.Args[4])
	case "-j":
		fmt.Println("Loading Local State...")
		states := "" // load from local persistent file
 		go dht.Want_to_join(ip, port, states)
		fmt.Printf("Stated Updated, Starting Join on [%s:%s]\n", ip, port)
	default:
		fmt.Println("The System Does NOT Understand This Option, Exiting...")
		os.Exit(-1)
	}


	transport.Fib(4)


	for {
		var op int
		fmt.Scanf("%d\n", &op)
		if op == 0 {
			os.Exit(-1)
		}
	}
}