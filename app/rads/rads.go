package main

import(
	"fmt"
)

func main() {
	fmt.Println("R advertising application started...")
	for {
		var s string
		fmt.Scanf("%s", &s)

		fmt.Println("selling ", s)
	}
}