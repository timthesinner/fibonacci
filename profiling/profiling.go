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
	heap := fibonacci.NewHeap(func(a, b int) int {
		return a - b
	})

	cpu, err := os.Create("HeapBulkCPU.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpu)
	defer func() {
		pprof.StopCPUProfile()
		cpu.Close()
	}()

	mem, err := os.Create("HeapBulkMEM.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		pprof.WriteHeapProfile(mem)
		mem.Close()
	}()

	for i := 0; i < count; i++ {
		heap.Insert(arr[i])
	}
	heap.Consolidate()
	fmt.Println("Finished Inserting")

	last := heap.RemoveMin()
	for i := 1; i < count; i++ {
		curr := heap.RemoveMin()
		if last > curr {
			fmt.Errorf("HeapBulkOperation(size=%d), %d > %d", heap.Size(), last, curr)
		}
		last = curr
	}

	fmt.Println(heap)
}
