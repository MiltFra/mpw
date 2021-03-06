package tree

import (
	"os"
	"sort"
	"sync"

	"github.com/miltfra/mpw/internal/out"
)

// leaf is a node without children
type leaf struct {
	T    *Tree  // the Tree this leaf belongs to
	p    *inner // the parent node
	keys [2]int // the keys (a <= x < b)
	f    string // the files that holds the map
}

// newleaf returns a new leaf with the given keys
func newLeaf(t *Tree, n *inner, keys ...int) *leaf {
	if len(keys) != 2 {
		out.Error("Invalid key count in leaf constructor.")
		return nil
	}
	p := dictPath(t, keys[0], keys[1])
	if _, err := os.Stat(p); os.IsNotExist(err) {
		writeToFile(make(map[int][]int), p)
	}
	return &leaf{t, n, [2]int{keys[0], keys[1]}, p}
}

// Insert adds values from a map into this leaf. If the key count
// exceeds the cap, the leaf is split into two and replaced by an inner
// node.
func (l *leaf) Insert(v map[int][]int, wg *sync.WaitGroup) node {
	defer wg.Done()
	contents := addMaps(v, readFromFile(l.f))
	if len(contents) > l.T.maxK {
		wg.Add(1)
		return l.Split(contents, wg)
	}
	writeToFile(contents, l.f)
	return l
}

func (l *leaf) SetKeys(keys [][]int) {
	if len(keys) != 1 || len(keys[0]) != 2 {
		out.Error("Only exactly 2 keys can be set to a leaf")
		return
	}

	l.keys = [2]int{keys[0][0], keys[0][1]}
}

// Split is called when the key count in the leaf exceeds a certain
// limit. Split then separates two equally big parts of the values
// and puts them into two different leaves. Then the leaf replaces
// itself with an inner node in the tree.
func (l *leaf) Split(v map[int][]int, wg *sync.WaitGroup) node {
	defer os.Remove(l.f)
	defer wg.Done()
	allKeys := make([]int, len(v))
	i := 0
	for k := range v {
		allKeys[i] = k
		i++
	}
	sort.Slice(allKeys, func(i, j int) bool {
		return allKeys[i] < allKeys[j]
	})
	keys := []int{
		l.keys[0],
		allKeys[len(allKeys)/2],
		l.keys[1],
	}
	node := newInner(l.T, keys...)
	// TODO:  Reuse this leaf object instead of allocating memory for a new one
	nodeL := newLeaf(l.T, node, keys[0], keys[1])
	nodeR := newLeaf(l.T, node, keys[1], keys[2])
	node.setChildren(nodeL, nodeR)
	wg.Add(1)
	go node.Insert(v, wg)
	if l.p != nil {
		if l.p.l == l {
			l.p.setChildren(node, l.p.r)
		} else {
			l.p.setChildren(l.p.l, node)
		}
	} else {
		l.T.root = node
	}
	return node
}

func (l *leaf) Tree() *Tree {
	return l.T
}

// Values returns the values in the leaf for a
// certain state
func (l *leaf) Values(state int) []int {
	dct := readFromFile(l.f)
	return dct[state]
}
