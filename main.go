package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
)

var (
	n      int
	coin   bool
	repeat uint
	lines  bool
)

var coins = []string{"HEADS", "tails"}

func init() {
	flag.IntVar(&n, "n", 0, "Defines [1, n] range; must be > 0")
	flag.BoolVar(&coin, "c", false, "Coin toss")
	flag.BoolVar(&lines, "i", false, "Input options as lines from stdin")
	flag.UintVar(&repeat, "r", 1, "Repeat count; must be > 0")

	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		fmt.Fprintln(flag.CommandLine.Output(), "\nNOTE: -n, -i, and -c are mutually exclusive.")
	}
}

func main() {
	flag.Parse()

	if repeat < 1 {
		fmt.Fprintln(os.Stderr, "-r argument must be > 0")
		flag.Usage()
		os.Exit(1)
	}

	var generateFunc func() any

	switch {
	case coin && n == 0 && !lines:
		generateFunc = func() any {
			r := rand.Intn(2)
			return coins[r]
		}

	case !coin && n != 0 && !lines:
		generateFunc = func() any {
			r := rand.Intn(n) + 1
			return r
		}

	case !coin && n == 0 && lines:
		options := []string{}
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			options = append(options, s.Text())
		}
		generateFunc = func() any {
			r := rand.Intn(len(options))
			return options[r]
		}

	default:
		flag.Usage()
		os.Exit(1)
	}

	results := []any{}
	for i := 0; i < int(repeat); i++ {
		results = append(results, generateFunc())
	}

	fmt.Println(results...)
}
