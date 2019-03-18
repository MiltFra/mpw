package main

import (
	"flag"
	"os"

	"github.com/miltfra/markov/internal/chain"
	"github.com/miltfra/markov/internal/out"
)

func main() {
	bufS := flag.Int("buf", 1024*256, "The buffer size to be used.")
	depth := flag.Int("n", 3, "The depth of the analysis?")
	file := flag.String("file", "", "Which file should be analyzed?")
	verbosity := flag.Int("verbosity", 3, "How verbose should the program be?")
	flag.Parse()
	if _, err := os.Stat(*file); os.IsNotExist(err) {
		out.Error("File does not exist...")
		return
	}
	out.Level = *verbosity
	chain.Analyze(*file, *depth, *bufS)
}
