package main

/*
	User Interfaces for both client & server side ops.
*/

import (
	"fmt"
	"os"
	"time"
	"github.com/rainer37/Rnet/dht"
	"github.com/rainer37/Rnet/transport"
)

func main() {
	fmt.Println("Initiating R-NET [", time.Now() ,"]")

	fmt.Println(os.Args[0])

	node := &dht.Plain_node{"localhost", 1388, "1338", "plain"}

	go dht.Self_init(node)

	transport.Fib(4)

	for {

	}
}