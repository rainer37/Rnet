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

	if len(os.Args) < 3 {
		fmt.Println("\aUsage:\n\tInit:\t-i [port] [topo]\n\tJoin:\t-j [target_ip] [target_port]")
		os.Exit(-1)
	}

	fmt.Println("Initiating R-NET v1.0.0\t[", time.Now() ,"]")
	fmt.Println("Local Arch & OS:\t[", runtime.GOARCH, ":", runtime.GOOS,"]")

	/* flags multiplexer
		-i : self initialization with [port] [topo]
		-j : join on an existing node if known ip with [target_ip] [target_port]
		-r : returning user
		TODO: MORE CASES -s -l -a
	*/
	
	flag, ip := os.Args[1], dht.Local_ip_4()

	switch flag {
	case "-i":
		port, topo := os.Args[2], os.Args[3] 
		go dht.Self_init(ip, port, topo)
	case "-j":
		port, tip := os.Args[3], os.Args[2] 
		states := "" // load from local persistent file
 		go dht.Want_to_join(tip, port, states)
		fmt.Printf("Starting Join on [%s:%s]\n", tip, port)
	case "-r":
		fmt.Println("Fetching Local States...")
		os.Exit(-1)
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