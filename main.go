package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"sync"
	"time"
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
	zero      bool
	comma     bool
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
	flag.BoolVar(&zero, "0", false, "\\0 delimiter in the output, similar to xargs -0 (shorthand)")
	flag.BoolVar(&zero, "null", false, "\\0 delimiter in the output, similar to xargs -0")
	flag.BoolVar(&comma, "d", false, "Comma delimiter in the output, similar to xargs -d, (shorthand)")

	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		fmt.Fprintln(flag.CommandLine.Output(), "\nNOTE: -n, -l, -t, and -c are mutually exclusive.")
		fmt.Fprintln(flag.CommandLine.Output(), "NOTE: -nl, -0, and -d are mutually exclusive.")
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
			r := rand.IntN(2)
			return coins[r]
		}

	case !coin && n != 0 && !lines && !tokens:
		generateFunc = func() any {
			r := rand.IntN(n) + 1
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
			r := rand.IntN(len(options))
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
			r := rand.IntN(len(options))
			return options[r]
		}

	default:
		flag.Usage()
		os.Exit(1)
	}

	// Produce and output results
	results := []any{}
	if shuffle {
		var wg sync.WaitGroup
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		wg.Add(1)
		go func(ctx context.Context, cancelFunc context.CancelFunc) {
			defer wg.Done()
			defer cancelFunc()

			// Results cannot repeat
			resultsKeys := make(map[any]interface{})
			for len(resultsKeys) < int(repeat) {
				r := generateFunc()
				if _, found := resultsKeys[r]; !found {
					results = append(results, r)
					resultsKeys[r] = nil
				}

				if ctx.Err() != nil {
					if ctx.Err().Error() == "context deadline exceeded" {
						fmt.Fprintln(os.Stderr, "Warning: timeout exceeded when generating results.  The repeat count may be too high.")
					}
					break
				}
			}
		}(ctx, cancelFunc)
		wg.Wait()

		rand.Shuffle(len(results), func(i, j int) {
			results[i], results[j] = results[j], results[i]
		})

	} else {
		for i := 0; i < int(repeat); i++ {
			results = append(results, generateFunc())
		}
	}

	switch {
	case newLine && !comma && !zero:
		for _, result := range results {
			fmt.Println(result)
		}
	case zero && !comma && !newLine:
		for _, result := range results {
			fmt.Printf("%v\x00", result)
		}
	case comma && !newLine && !zero:
		for _, result := range results {
			fmt.Printf("%v,", result)
		}

	case !comma && !newLine && !zero:
		fmt.Println(results...)

	default:
		flag.Usage()
		os.Exit(1)
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
