package chain

import "testing"

func TestAnalyze(t *testing.T) {
	Analyze("/home/miltfra/go/src/github.com/miltfra/markov/test/analyze_test/example1.txt", 3, 1024)
}
