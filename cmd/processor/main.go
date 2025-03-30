package main

import (
	"flag"
	"log"
	"sync"
)

var (
	goroutines = flag.Int("workers", 100, "Initial number of goroutines to create during startup")
)

func main() {
	flag.Parse()
	run()
}
