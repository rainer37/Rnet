package main

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

	dht.Echo("Sup!")
	transport.Fib(4)
}