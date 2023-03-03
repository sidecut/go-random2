package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

var (
	n         int
	coin      bool
	repeat    uint
	repeatSet bool
	lines     bool
	tokens    bool
	shuffle   bool
	newLine   bool
)

const (
	shuffleUsage = "Shuffle the input"
)

var coins = []string{"HEADS", "tails"}

func init() {
	flag.IntVar(&n, "n", 0, "Defines [1, n] range; must be > 0")
	flag.BoolVar(&coin, "c", false, "Coin toss")
	flag.BoolVar(&lines, "l", false, "Input options as lines from stdin")
	flag.BoolVar(&tokens, "t", false, "Input options as tokens from stdin")
	flag.UintVar(&repeat, "r", 1, "Repeat count; must be > 0")
	flag.BoolVar(&shuffle, "shuffle", false, shuffleUsage)
	flag.BoolVar(&shuffle, "s", false, shuffleUsage+" (shorthand)")
	flag.BoolVar(&newLine, "nl", false, "Newline between items in the output")

	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		fmt.Fprintln(flag.CommandLine.Output(), "\nNOTE: -n, -l, -t, and -c are mutually exclusive.")
	}
}

func main() {
	// Parse and validate options
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "r" {
			repeatSet = true
		}
	})

	// If we're shuffling, then repeat defaults to the number of items in the input set
	if repeat < 1 && !shuffle {
		fmt.Fprintln(os.Stderr, "-r argument must be > 0")
		flag.Usage()
		os.Exit(1)
	}

	// Create appropriate generator
	var generateFunc func() any

	switch {
	case coin && n == 0 && !lines && !tokens:
		generateFunc = func() any {
			r := rand.Intn(2)
			return coins[r]
		}

	case !coin && n != 0 && !lines && !tokens:
		generateFunc = func() any {
			r := rand.Intn(n) + 1
			return r
		}

	case !coin && n == 0 && lines && !tokens:
		options := []string{}
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			options = append(options, s.Text())
		}
		if len(options) == 0 {
			fmt.Fprintln(os.Stderr, "Warning: no options specified")
		}
		if len(options) == 1 {
			fmt.Fprintln(os.Stderr, "Warning: only one option specified.  Did you mean -t?")
		}
		if shuffle {
			setOrValidateRepeat(options)
		}

		generateFunc = func() any {
			r := rand.Intn(len(options))
			return options[r]
		}

	case !coin && n == 0 && !lines && tokens:
		options := []string{}
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			fields := strings.Fields(s.Text())
			options = append(options, fields...)
		}
		if shuffle {
			setOrValidateRepeat(options)
		}

		generateFunc = func() any {
			r := rand.Intn(len(options))
			return options[r]
		}

	default:
		flag.Usage()
		os.Exit(1)
	}

	// Produce and output results
	results := []any{}
	if shuffle {
		// Results cannot repeat
		resultsKeys := make(map[any]interface{})
		for len(resultsKeys) < int(repeat) {
			r := generateFunc()
			if _, found := resultsKeys[r]; !found {
				results = append(results, r)
				resultsKeys[r] = nil
			}
		}
		rand.Shuffle(len(results), func(i, j int) {
			results[i], results[j] = results[j], results[i]
		})

	} else {
		for i := 0; i < int(repeat); i++ {
			results = append(results, generateFunc())
		}
	}

	if newLine {
		for _, result := range results {
			fmt.Println(result)
		}
	} else {
		fmt.Println(results...)
	}
}

func setOrValidateRepeat(options []string) {
	if !repeatSet {
		repeat = uint(len(options))
	} else if repeat > uint(len(options)) {
		fmt.Fprintln(os.Stderr, "Repeat count cannot exceed the number of options")
		os.Exit(1)
	}
}
