package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
	err := exec.Command("resolvectl", "flush-caches").Run()
	if err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	lookup := func(n int) {
		defer wg.Done()
		out, err := exec.Command("resolvectl", "query", "manual-inject.default").CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%d: %v", n, err.Error())
		} else {
			fmt.Printf("%d IP %v\n", n, string(out))
		}
	}
	go lookup(1)
	time.Sleep(300 * time.Millisecond)
	go lookup(2)
	wg.Wait()
}
