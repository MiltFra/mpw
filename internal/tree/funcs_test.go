package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDictPath(t *testing.T) {
	tr := &Tree{dir: "path"}
	if dictPath(tr, 0, 1) != "path/0_1" {
		t.FailNow()
	}
}

func TestReadWrite(t *testing.T) {
	p := "/home/miltfra/test/0_1"
	dct := map[int][]int{0: []int{1, 2, 3}}
	writeToFile(dct, p)
	dct2 := readFromFile(p)
	for k := range dct {
		for i := range dct[k] {
			if dct[k][i] != dct2[k][i] {
				t.FailNow()
			}
		}
	}
}

func TestAddSlices(t *testing.T) {
	arr1 := []int{1, 2, 3, 4}
	arr2 := []int{1, 2, 3, 4}
	arr3 := addSlices(arr1, arr2)
	for i := range arr1 {
		if arr1[i]+arr2[i] != arr3[i] {
			t.FailNow()
		}
	}
}

func TestAddMaps(t *testing.T) {
	dct1 := map[int][]int{0: []int{1, 2, 3}}
	dct2 := map[int][]int{0: []int{1, 2, 3}, 1: []int{1, 2, 3}}
	dct3 := addMaps(dct1, dct2)
	for k := range dct3 {
		arr := addSlices(dct1[k], dct2[k])
		assert.Equal(t, dct3[k], arr)
	}
}

func TestSplitMap(t *testing.T) {
	dct := map[int][]int{0: []int{1, 2, 3}, 1: []int{1, 2, 3}, 2: []int{0, 4}}
	key := 1
	dct1, dct2 := splitMap(key, dct)
	assert.Equal(t, addMaps(dct1, dct2), dct)
}
