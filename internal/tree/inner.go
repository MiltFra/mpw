package tree

import (
	"sync"

	"github.com/miltfra/markov/internal/out"
)

// innerNode is a part of a Tree that is not a Leaf.
type inner struct {
	T    *Tree  // the parent
	keys [3]int // the keys of this node; a, b, c
	l    node   // holds a <= x < b
	r    node   // holds b <= x < c
}

// newNode returns a new node with given keys
func newInner(t *Tree, keys ...int) *inner {
	if len(keys) != 3 {
		out.Error("Invalid key count in innerNode constructor.")
		return nil
	}
	return &inner{t, [3]int{keys[0], keys[1], keys[2]}, nil, nil}
}

func (in *inner) setChildren(l, r node) {
	in.l = l
	in.r = r
}

func (in *inner) Insert(v map[int][]int, wg *sync.WaitGroup) node {
	defer wg.Done()
	x, y := splitMap(in.keys[1], v)
	wg.Add(2)
	go in.l.Insert(x, wg)
	go in.r.Insert(y, wg)
	return in
}

func (in *inner) Values(state int) []int {
	if state < in.keys[1] {
		return in.l.Values(state)
	}
	return in.r.Values(state)
}

func (in *inner) Tree() *Tree {
	return in.T
}
