package tree

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/miltfra/markov/internal/out"
)

// dictPath returns a path corresponding to a given pair of keys
func dictPath(t *Tree, a, b int) string {
	return fmt.Sprintf("%v/%v_%v", t.dir, a, b)
}

// writeToFile writes a map to a given file
func writeToFile(v map[int][]int, p string) {
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
	err = encoder.Encode(v)
	if err != nil {
		out.Error("Could not encode a map to a file.")
		panic(err)
	}
	out.Status("Successfully wrote to (" + p + ")")
}

// readFromFile reads a map from a given file
func readFromFile(p string) map[int][]int {
	r, err := os.Open(p)
	defer r.Close()
	if err != nil {
		out.Error("Could not read from a file. (" + p + ")")
		panic(err)
	}
	decoder := gob.NewDecoder(r)
	var v map[int][]int
	err = decoder.Decode(&v)
	if err != nil {
		out.Error("Could not decode file.")
		panic(err)
	}
	return v
}

func addMaps(x, y map[int][]int) (z map[int][]int) {
	z = make(map[int][]int)
	for k := range x {
		z[k] = x[k]
	}
	for k, v := range y {
		z[k] = addSlices(z[k], v)
	}
	return
}

func addSlices(x, y []int) (z []int) {
	l := len(x)
	if l < len(y) {
		l = len(y)
	}
	z = make([]int, l)
	for i := 0; i < l; i++ {
		if i < len(x) {
			z[i] += x[i]
		}
		if i < len(y) {
			z[i] += y[i]
		}
	}
	return
}

// splitMap returns two parts of a map where the first contains
// all the key-value pairs where the key is less than the given
// one.
func splitMap(key int, z map[int][]int) (map[int][]int, map[int][]int) {
	x := make(map[int][]int)
	y := make(map[int][]int)
	for k, v := range z {
		if k < key {
			x[k] = v
		} else {
			y[k] = v
		}
	}
	return x, y
}
