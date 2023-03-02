package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
)

var (
	n      int
	coin   bool
	repeat uint
)

func init() {
	flag.IntVar(&n, "n", 0, "Defines [1, n] range; must be > 0")
	flag.BoolVar(&coin, "c", false, "Coin toss")
	flag.UintVar(&repeat, "r", 1, "Repeat count; must be > 0")

	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		fmt.Fprintln(flag.CommandLine.Output(), "\nNOTE: -n and -c are mutually exclusive.")
	}
}

func main() {
	flag.Parse()

	if repeat < 1 {
		fmt.Fprintln(os.Stderr, "-r argument must be > 0")
		flag.Usage()
		os.Exit(1)
	}

	switch {
	case coin && n == 0:
		for i := 0; i < int(repeat); i++ {
			coins := []string{"HEADS", "tails"}
			r := rand.Intn(2)
			fmt.Println(coins[r])
		}

	case !coin && n != 0:
		for i := 0; i < int(repeat); i++ {
			r := rand.Intn(n-1) + 1
			fmt.Println(r)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}
