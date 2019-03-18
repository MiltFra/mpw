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

func (in *inner) SetKeys(keys [][]int) {
	if len(keys) < 2 {
		out.Error("At least two pairs of keys need to be set to a leaf.")
		return
	}
	center := len(keys) / 2
	in.keys = [3]int{keys[0][0], keys[center][0], keys[len(keys)-1][1]}
	var l, r node
	if center == 1 {
		l = newLeaf(in.T, in, keys[0]...)
	} else {
		l = newInner(in.T, 0, 0, 0)
		l.SetKeys(keys[:center])
	}
	if len(keys) == 2 {
		r = newLeaf(in.T, in, keys[1]...)
	} else {
		r = newInner(in.T, 0, 0, 0)
		r.SetKeys(keys[center:])
	}
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
