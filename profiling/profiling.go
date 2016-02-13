// profiling.go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"

	"github.com/timthesinner/fibonacci"
)

func main() {
	const count = 10000000
	arr := [count]int{}
	for i := 0; i < count; i++ {
		arr[i] = rand.Int()
	}
	heap := fibonacci.NewHeap(func(a, b interface{}) int {
		return a.(int) - b.(int)
	})

	f, err := os.Create("HeapBulkCPU.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for i := 0; i < count; i++ {
		heap.Insert(arr[i])
	}

	last := heap.RemoveMin().(int)
	for i := 1; i < count; i++ {
		curr := heap.RemoveMin().(int)
		if last > curr {
			fmt.Errorf("HeapBulkOperation(size=%d), %d > %d", heap.Size(), last, curr)
		}
		last = curr
	}

	fmt.Println(heap)
}
