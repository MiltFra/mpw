package tree

import (
	"sync"

	"github.com/miltfra/markov/internal"
)

// Tree manages partial maps on the drive to
// read them when necessary
type Tree struct {
	root node   // the first node
	maxK int    // maximum number of keys per leaf
	dir  string // the path of the directory of the Tree
}

type node interface {
	Insert(map[int][]int, *sync.WaitGroup) node
	Values(int) []int
	Tree() *Tree
}

// New returns a pointer to a new Tree object with
// a depth and a directory
func New(n int, d string) *Tree {
	internal.ResetDir(d)
	a := 0
	b := internal.GetPower(n) + 1
	t := &Tree{nil, 100000, d}
	r := newLeaf(t, nil, a, b)
	t.root = r
	return t
}

// Insert inserts a map into the tree.
// This function makes heavy use of concurrency.
func (t *Tree) Insert(v map[int][]int) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	t.root.Insert(v, wg)
	wg.Wait()
}
