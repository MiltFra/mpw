package main

import (
	"os"

	"github.com/miltfra/markov/internal/chain"
	"github.com/miltfra/markov/internal/out"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		out.Error("Unexpected number of arguments; aborting...")
	}
	chain.Analyze(args[0], 4)
}
