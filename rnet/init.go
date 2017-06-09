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
	
	fmt.Println(os.Args[0])

	node := &dht.Plain_node{"localhost", 1388, "1338", "plain"}

	go dht.Self_init(node)

	transport.Fib(4)

	
	for {
		var op int
		fmt.Scanf("%d\n", &op)
		if op == 0 {
			os.Exit(-1)
		}
	}
}