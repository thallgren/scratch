package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	ls, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP: net.IPv4(172,17,0,1),
		Port: 34567,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ls.Close()
	fmt.Println("success")
}