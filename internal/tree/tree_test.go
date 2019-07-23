package tree

import (
	"runtime"
	"testing"

	"github.com/miltfra/mpw/internal"
)

func TestInsert(t *testing.T) {
	runtime.GOMAXPROCS(8)
	tr := New(70, "/home/miltfra/go/src/github.com/miltfra/mpw/test/testinsert-mc")
	p := internal.GetPower(2)
	for i := 0; i < 50; i++ {
		values := make(map[int][]int)
		for j := 0; j < p; j++ {
			values[i*p+j] = make([]int, 95)
			values[i*p+j][0]++
		}
		tr.Insert(values)
	}
}
