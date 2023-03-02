package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
)

var (
	n int
)

func init() {
	flag.IntVar(&n, "n", 0, "Defines [1, n] range.\nn must be > 0.")
}

func main() {
	flag.Parse()
	if n < 1 {
		flag.Usage()
		os.Exit(1)
	}
	r := rand.Intn(n-1) + 1
	fmt.Printf("%v\n", r)
}
