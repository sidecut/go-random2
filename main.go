package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
)

var (
	n    int
	coin bool
)

func init() {
	flag.IntVar(&n, "n", 0, "Defines [1, n] range.\nn must be > 0.")
	flag.BoolVar(&coin, "c", false, "Coin toss.")
	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		fmt.Fprintln(flag.CommandLine.Output(), "\nNOTE: -n and -c are mutually exclusive.")
	}
}

func main() {
	flag.Parse()

	switch {
	case coin && n == 0:
		coins := []string{"HEADS", "tails"}
		r := rand.Intn(2)
		fmt.Println(coins[r])

	case !coin && n != 0:
		r := rand.Intn(n-1) + 1
		fmt.Println(r)
	default:
		flag.Usage()
		os.Exit(1)
	}
}
