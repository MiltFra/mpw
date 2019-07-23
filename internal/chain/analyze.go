package chain

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/miltfra/mpw/internal"
	"github.com/miltfra/mpw/internal/out"
)

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

func update(b *int, size float64) {
	for *b != -1 {
		fmt.Printf("Read %v bytes (%v%v)        \r", *b, float64(*b)*100/size, "%")
		time.Sleep(1 * time.Second)
	}
}
