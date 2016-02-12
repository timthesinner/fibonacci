// fibonacci_test
package fibonacci

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"testing"
)
import "math/rand"

func TestHeapNode(t *testing.T) {
	cases := []struct{ in, want interface{} }{
		{"test", "test"},
		{1, 1},
	}
	for _, c := range cases {
		got := HeapNode{value: c.in}
		if got.value != c.want {
			t.Errorf("HeapNode(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestHeapCreate(t *testing.T) {
	heap := NewHeap(func(a, b interface{}) int {
		return a.(int) - b.(int)
	})

	cases := []struct{ in, want int }{
		{100, 100},
		{101, 100},
		{99, 99},
	}
	for _, c := range cases {
		heap.insert(c.in)
		if heap.peek() != c.want {
			t.Errorf("HeapInsert(%q) == %q, wanted %q", c.in, heap.peek(), c.want)
		}
	}

	for i := 10; i >= 0; i-- {
		heap.insert(i)
		if heap.peek() != i {
			t.Errorf("HeapInsert(%q) == %q, wanted %q", i, heap.peek(), i)
		}
	}

	for i := 0; i < 14; i++ {
		removed := heap.removeMin()
		expected := i
		if i > 10 {
			expected = i - 11 + 99
		}
		if removed != expected {
			t.Errorf("HeapRemoveMin == %q, wanted %d", removed, expected)
		}
	}
}

func TestHeapBulk(t *testing.T) {
	heap := NewHeap(func(a, b interface{}) int {
		return a.(int) - b.(int)
	})

	count := 10000
	for i := 0; i < count; i++ {
		heap.insert(rand.Int())
	}

	last := heap.removeMin().(int)
	for i := 1; i < count; i++ {
		curr := heap.removeMin().(int)
		if last > curr {
			t.Errorf("HeapBulkOperation(size=%d), %d > %d", heap.size, last, curr)
		}
		last = curr
	}

	fmt.Println(heap)
}
