package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	a, err := net.ResolveUDPAddr("udp", "127.0.0.1:")
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenUDP("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	_ = l.SetReadDeadline(time.Now().Add(time.Second))
	buf := [0x10000]byte{}
	n, _, err := l.ReadFrom(buf[:])
	if n > 0 {
		fmt.Println(string(buf[:n]))
	}
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			log.Printf("timeout = %t", opErr.Timeout())
		}
		log.Fatal(err)
	}
}
