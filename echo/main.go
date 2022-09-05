package main

import (
	"os"
)

func main() {
	for _, s := range os.Args[1:] {
		os.Stdout.Write([]byte{'"'})
		os.Stdout.Write([]byte(s))
		os.Stdout.Write([]byte{'"', '\n'})
	}
}
