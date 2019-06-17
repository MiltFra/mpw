package chain

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/miltfra/markov/internal"
)

// Generate loads a markov chain and uses the probability
// distribution to generate a number of strings in a
// certain size intervall. To not interfere with the
// randomness strings are discarded until the match the
// length that was given. Thus it might take a while to
// get a big number of fitting strings.
func Generate(meta string, l0, l1, count int) {
	c := readMeta(meta)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < count; i++ {
		fmt.Println(c.randomWord(l0, l1))
	}
}

func (c *Chain) randomWord(l0, l1 int) string {
	cha := c.randomCandidate()
	for len(cha) < l0 || len(cha) >= l1 {
		cha = c.randomCandidate()
	}
	return arrToString(cha)
}

func arrToString(arr []int) string {
	s := ""
	for _, v := range arr {
		s += string(v)
	}
	return s
}

func (c *Chain) randomCandidate() []int {
	// TODO: Reuse the byte array and discard the candidates in the same function.
	candidate := make([]int, 0, 20)
	state := 0
	cha := c.randomSymbol(state)
	state = internal.ExtendState(c.N, state, cha)
	candidate = append(candidate, cha+31)
	for state != 0 {
		cha = c.randomSymbol(state)
		state = internal.ExtendState(c.N, state, cha)
		candidate = append(candidate, cha+31)
	}
	return candidate
}

func (c *Chain) randomSymbol(state int) int {
	return weightedCoinflip(c.values.Values(state))
}

func weightedCoinflip(values []int) int {
	// get the sum
	s := 0
	for _, v := range values {
		s += v
	}
	// if there haven't been any occurences we return 0
	if s == 0 {
		return 0
	}
	r := rand.Intn(s)
	for i, v := range values {
		r -= v
		if r < 0 {
			return i
		}
	}
	return 0
}
