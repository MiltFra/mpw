package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/miltfra/markov/internal/chain"
	"github.com/miltfra/markov/internal/out"
)

var (
	flagsAnalyze  = flag.NewFlagSet("analyze", flag.ContinueOnError)
	flagsGenerate = flag.NewFlagSet("generate", flag.ContinueOnError)
	verbosity     int
	cores         int
	buffer        int
	n             int
	l0            int
	l1            int
	count         int
)

func main() {
	flagsAnalyze.IntVar(&verbosity, "silence", 2, "How silent should the program be?")
	flagsAnalyze.IntVar(&cores, "cores", runtime.NumCPU(), "How many cores should the program use?")
	flagsAnalyze.IntVar(&buffer, "buffer", 256*1024, "How many bytes should be read at once?")
	flagsAnalyze.IntVar(&n, "n", 3, "How many previous characters should be taken into account?")

	flagsGenerate.IntVar(&verbosity, "silence", 2, "How silent should the program be?")
	flagsGenerate.IntVar(&cores, "cores", runtime.NumCPU(), "How many cores should the program use?")
	flagsGenerate.IntVar(&l0, "l", 0, "What is the MINIMUM length of the desired words?")
	flagsGenerate.IntVar(&l1, "L", 1<<31, "What is the MAXIMUM length of the desired words?")
	switch os.Args[1] {
	case "analyze":
		analyze()
	case "generate":
		generate()
	}
}

func analyze() {
	flagsAnalyze.Parse(os.Args[2:])
	file := flagsAnalyze.Arg(0)
	// Running the program
	// Checking for errors
	if _, err := os.Stat(file); os.IsNotExist(err) {
		out.Error(fmt.Sprintf("File does not exist...(%v)", file))
		return
	}
	// Setting config
	out.Level = verbosity
	runtime.GOMAXPROCS(cores)
	// Analyzing
	chain.Analyze(file, n, buffer)
}

func generate() {
	flagsAnalyze.Parse(os.Args[2:])
	file := flagsAnalyze.Arg(0)
	// Running the program
	// Checking for errors
	if _, err := os.Stat(file); os.IsNotExist(err) {
		out.Error(fmt.Sprintf("File does not exist...(%v)", file))
		return
	}
	// Setting config
	out.Level = verbosity
	runtime.GOMAXPROCS(cores)
	chain.Generate(file, l0, l1, count)
}
