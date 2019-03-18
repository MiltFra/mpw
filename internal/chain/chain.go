package chain

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/miltfra/markov/internal"
	"github.com/miltfra/markov/internal/out"
	"github.com/miltfra/markov/internal/tree"
)

// Chain contains probabilites of transitions
// from one state to another. Technically the
// occurences of each transition are counted.
type Chain struct {
	N      int
	path   string
	values *tree.Tree
	final  bool
}

// Insert adds a given map of values into the chain
// if it has not yet been finalized. Otherwise it will
// print a warning and ignore the insert.
func (c *Chain) Insert(v map[int][]int) {
	if c.final {
		out.Warning("Cannot insert into finalized chain.")
		return
	}
	c.values.Insert(v)
}

// New returns a Pointer to a new Chain object based on a
// directory d with a depth n.
func New(n int, d string) *Chain {
	return &Chain{n, d, tree.New(n, d), false}
}

// Analyze reads a file and returns a chain which contains
// probabilities corresponding to n characters followed by
// another one.
func Analyze(path string, n int) (c *Chain) {
	c = New(n, path+"-mc")
	file, err := os.Open(path)
	if err != nil {
		out.Error("Could not open input. (" + path + ")")
		panic(err)
	}
	defer file.Close()
	buf := make([]byte, 1024*256)
	state := 0
	read := 0
	go update(&read)
	var dct map[int][]int
	for {
		dct = make(map[int][]int)
		count, err := file.Read(buf)
		if count > 0 {
			read += count
			out.Status(fmt.Sprintf("Read %v bytes; processing...", count))
			actBuf := buf[:count]
			for _, s := range actBuf {
				s -= 31
				if s > 95 {
					s = 0
				}
				if dct[state] == nil {
					dct[state] = make([]int, 96)
				}
				dct[state][s]++
				state = internal.ExtendState(n, state, int(s))
			}
			c.Insert(dct)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			out.Error("Could not read input. (" + path + ")")
			break
		}
	}
	read = -1
	out.Status("Reading file complete.")
	return
}

func update(b *int) {
	for *b != -1 {
		out.Status(fmt.Sprintf("[STA] Read %v bytes\r", *b))
		time.Sleep(1 * time.Second)
	}
}

// Finalize converts the occurences into probabilities thus
// making the chain unable to be edited.
func (c *Chain) Finalize() {

}
