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
	ac "github.com/rainer37/Rnet/acontroller"
)

func load_local_states() {

}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("\aUsage:\n\tInit:\t-i [port] [topo]\n\tJoin:\t-j [target_ip] [target_port] [my_port]")
		os.Exit(-1)
	}

	loc_ip := dht.Local_ip_4()

	fmt.Println("Initiating R-NET v1.0.0\t[", time.Now() ,"]")
	fmt.Println("Local Arch & OS:\t[", runtime.GOARCH, ":", runtime.GOOS,"]")
	fmt.Printf("PID:\t\t\t[%d]\n", os.Getpid())
	fmt.Printf("Local IP:\t\t[%s]\n", loc_ip)
	/* flags multiplexer
		-i : self initialization with [port] [topo]
		-j : join on an existing node if known ip with [target_ip] [target_port] [my_port]
		-r : returning user
		TODO: MORE CASES -s -l -a
	*/

	var por string
	flag, ip := os.Args[1], loc_ip

	switch flag {
	case "-i":
		port, topo := os.Args[2], os.Args[3]
		por = port 
		go dht.Self_init(ip, port, topo)
	case "-j":
		port, tip, my_port := os.Args[3], os.Args[2], os.Args[4]
		states := "" // load from local persistent file
 		por = my_port
 		go dht.Want_to_join(tip, port, states, my_port)
		fmt.Printf("Starting Join on [%s:%s]\n", tip, port)
	case "-r":
		fmt.Println("Fetching Local States...")
		os.Exit(-1)
	default:
		fmt.Println("The System Does NOT Understand This Option, Exiting...")
		os.Exit(-1)
	}

	time.Sleep(1 * time.Second) // wait for 1s for nice fmt.
	
	go transport.Trans_Boot(por)
	// booting the App controller
	ac.AC_boot()
	// simple client shell
	for {
		fmt.Print("Rnet:r$ ")
		var op string
		fmt.Scanf("%s\n", &op)

		switch op {
		case "exit","0","q","quit":
			transport.Close_ln()
			os.Exit(-1)
		case "ls":
			ac.Get_app_list()
		case "dht":
			dht.Print()
		case "":
			fmt.Println()
		case "1":
			ac.Exec_app("chatchat/chat") // sample application 1.
		default:
			fmt.Println("Unknown CMD.")
		}
	}
}
