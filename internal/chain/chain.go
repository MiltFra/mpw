package chain

import (
	"encoding/gob"
	"os"
	"strings"

	"github.com/miltfra/mpw/internal/out"
	"github.com/miltfra/mpw/internal/tree"
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

type meta struct {
	N int
}

func (c *Chain) writeMeta() {
	p := c.path + "/meta"
	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.Create(p)
	}
	w, err := os.OpenFile(p, os.O_WRONLY, os.ModeExclusive)
	defer w.Close()
	if err != nil {
		out.Error("Could not write to a file. (" + p + ")")
		panic(err)
	}
	encoder := gob.NewEncoder(w)
	err = encoder.Encode(meta{c.N})
	if err != nil {
		out.Error("Could not encode a map to a file.")
		panic(err)
	}
	out.Status("Successfully wrote to (" + p + ")")
}

func readMeta(p string) *Chain {
	r, err := os.Open(p)
	defer r.Close()
	if err != nil {
		out.Error("Could not read from a file. (" + p + ")")
		panic(err)
	}
	decoder := gob.NewDecoder(r)
	var m meta
	err = decoder.Decode(&m)
	if err != nil {
		out.Error("Could not decode file.")
		panic(err)
	}
	split := strings.Split(p, "/")
	d := strings.Join(split[:len(split)-1], "/")
	return &Chain{m.N, d, tree.Load(m.N, d), false}
}
