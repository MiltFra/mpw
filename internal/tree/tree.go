package tree

import (
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/miltfra/markov/internal/out"

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
	SetKeys([][]int)
}

var defaultMaxK = 100000

// New returns a pointer to a new Tree object with
// a depth and a directory
func New(n int, d string) *Tree {
	internal.ResetDir(d)
	a := 0
	b := internal.GetPower(n) + 1
	t := &Tree{nil, defaultMaxK, d}
	r := newLeaf(t, nil, a, b)
	t.root = r
	return t
}

// Load returns a pointer to a new Tree object base
// on a directory structure.
func Load(n int, d string) *Tree {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		out.Error("Could not read directory. (" + d + ")")
	}
	keys := make([][]int, 0, len(files))
	var k []int
	for _, f := range files {
		k = parseFile(f.Name())
		if k != nil {
			keys = append(keys, k)
		}
	}
	t := &Tree{nil, defaultMaxK, d}
	var root node
	if len(keys) == 1 {
		b := internal.GetPower(n) + 1
		root = newLeaf(t, nil, 0, b)
	} else {
		root = newInner(t, 0, 0, 0)
		root.SetKeys(keys)
	}
	t.root = root
	return t
}

func parseFile(f string) []int {
	s := strings.Split(f, "_")
	if len(s) != 2 {
		return nil
	}
	i1, err := strconv.Atoi(s[0])
	if err != nil {
		return nil
	}
	i2, err := strconv.Atoi(s[1])
	if err != nil {
		return nil
	}
	return []int{i1, i2}
}

// Insert inserts a map into the tree.
// This function makes heavy use of concurrency.
func (t *Tree) Insert(v map[int][]int) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	t.root.Insert(v, wg)
	wg.Wait()
}

// Values returns the slice with the values of a certain
// state.
func (t *Tree) Values(state int) []int {
	return t.root.Values(state)
}
