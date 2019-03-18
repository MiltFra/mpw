package chain

import "testing"

func TestNew(t *testing.T) {
	New(5, "/home/miltfra/go/src/github.com/miltfra/markov/test/testnew")
}

func TestAnalyze(t *testing.T) {
	Analyze("/home/miltfra/go/src/github.com/miltfra/markov/test/gamigo_10_11", 3)
}
