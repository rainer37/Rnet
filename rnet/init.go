package main

/*
	User Interfaces for both client & server side ops.

	/* flags multiplexer
	-i : self initialization with [port] [topo]
	-j : join on an existing node if known ip with [target_ip] [target_port] [my_port]
	-r : returning user
	TODO: MORE CASES -s -l -a
*/

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"bufio"
	"strings"
	"github.com/rainer37/Rnet/dht"
	ac "github.com/rainer37/Rnet/apctl"
)

func load_local_states() {

}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("\aUsage:\n\tInit:\t-i [port] [topo]\n\tJoin:\t-j [target_ip] [target_port] [my_port]")
		os.Exit(-1)
	}

	loc_ip := dht.Local_ip_4()
	ARCH, OS := runtime.GOARCH, runtime.GOOS

	fmt.Printf("Initiating R-NET v1.0.0\t[ %s ]\n", time.Now())
	fmt.Printf("Local Arch & OS:\t[ %v : %v ]\n", ARCH, OS)
	fmt.Printf("Current Rnet PID:\t[ %d ]\n", os.Getpid())
	fmt.Printf("Local IP:\t\t[ %s ]\n", loc_ip)

	flag, ip := os.Args[1], loc_ip

	switch flag {
	case "-i":
		port, topo := os.Args[2], os.Args[3]
		go dht.Self_init(ip, port, topo)
	case "-j":
		port, tip, my_port := os.Args[3], os.Args[2], os.Args[4]
		states := "" // load from local persistent file
 		go dht.Want_to_join(tip, port, states, my_port)
		fmt.Printf("Starting Join on [%s:%s]\n", tip, port)
	case "-r":
		fmt.Println("Fetching Local States...")
		os.Exit(-1)
	default:
		fmt.Println("The System Does NOT Understand This Option, Exiting...")
		os.Exit(-1)
	}

	
	// booting the App controller
	go ac.AC_boot()

	time.Sleep(1 * time.Second) // wait for 1s for nice fmt.


	// simple client shell
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Rnet:r$ ")

	for scanner.Scan() {
		var text string = scanner.Text()
		cmd := strings.Split(text, " ")
		op := cmd[0]

		switch op {
		case "exit","0","q","quit":
			ac.Close_apc_uds_ln()
			os.Exit(-1)
		case "ls":
			ac.Get_app_list()
		case "dht":
			dht.Print()
		case "m":
			// mails
		case "":
			fmt.Println()
		case "e": // exec application
			if len(cmd) > 1 {
				app_name := cmd[1]
				fmt.Println("Executing "+app_name)
				err := ac.Exec_app(app_name, OS) // sample application 1.
				if err != nil {
					fmt.Println("Cannot execute application"+app_name)
				}
			} else {
				fmt.Println("No application specified. Usage: e [app_name]")
			}
		default:
			fmt.Println("Unknown CMD.")
		}

		fmt.Print("Rnet:r$ ")
	}
}
