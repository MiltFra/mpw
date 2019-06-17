package chain

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strings"
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
func Analyze(path string, n int, bufS int) (c *Chain) {
	c = New(n, path+"-mc")
	file, err := os.Open(path)
	if err != nil {
		out.Error("Could not open input. (" + path + ")")
		panic(err)
	}
	defer file.Close()
	buf := make([]byte, bufS)
	state := 0
	read := 0
	c.writeMeta()
	f, err := os.Stat(path)
	go update(&read, float64(f.Size()))
	var dct map[int][]int
	for {
		dct = make(map[int][]int)
		count, err := file.Read(buf)
		if count > 0 {
			read += count
			//out.Status(fmt.Sprintf("Read %v bytes; processing...", count))
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
	// TODO: Write meta to disk
	return
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

func update(b *int, size float64) {
	for *b != -1 {
		fmt.Printf("Read %v bytes (%v%v)        \r", *b, float64(*b)*100/size, "%")
		time.Sleep(1 * time.Second)
	}
}
