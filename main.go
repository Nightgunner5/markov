package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	seed   = flag.Int64("seed", time.Now().UnixNano(), "the random seed to use when generating text.")
	output = flag.String("o", "", "if non-empty, the file to write output to. defaults to standard output.")
	input  = flag.String("i", "", "if non-empty, the file to read input from. defaults to standard input.")
	c      = flag.Bool("c", false, "compile the input to a reusable markov chain cache instead of generating text.")
	p      = flag.Bool("p", false, "if this is set, the input will be parsed as a reusable markov chain cache.")
	n      = flag.Int("n", 5, "the number of lines of text to generate")
)

func main() {
	var err error
	flag.Parse()

	r := rand.New(rand.NewSource(*seed))

	in := os.Stdin
	if *input != "" {
		in, err = os.Open(*input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while opening -i file: %s\n", err)
			return
		}
		defer in.Close()
	}

	out := os.Stdout
	if *output != "" {
		out, err = os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while opening -o file: %s\n", err)
			return
		}
		defer out.Close()
	}

	chain := NewChain(2)

	if *p {
		g, err := gzip.NewReader(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while reading precompiled file: %s\n", err)
			return
		}

		err = json.NewDecoder(g).Decode(&chain)
		g.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "error while reading precompiled file: %s\n", err)
			return
		}
	} else {
		b := bufio.NewReader(in)

		for {
			line, err := b.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "error while reading input: %s\n", err)
				return
			}

			chain.Build(strings.Fields(line))
		}
	}

	if *c {
		g, _ := gzip.NewWriterLevel(out, gzip.BestCompression)
		err = json.NewEncoder(g).Encode(chain)
		g.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "error while writing compiled file: %s\n", err)
			return
		}
	} else {
		for i := 0; i < *n; i++ {
			fmt.Fprintln(out, chain.Generate(100, r))
		}
	}
}
