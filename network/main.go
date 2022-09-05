package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(addr.String())
}
