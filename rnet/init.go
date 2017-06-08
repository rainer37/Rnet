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

	dht.Self_init("192.168.0.31", 1388, "chord")

	transport.Fib(4)
}