// fibonacci_test
package fibonacci

import (
	"fmt"
	"math/rand"
	"testing"
)

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
		{5577006791947779410, 99},
		{8674665223082153551, 99},
		{6129484611666145821, 99},
		{4037200794235010051, 99},
		{3916589616287113937, 99},
		{6334824724549167320, 99},
		{605394647632969758, 99},
		{1443635317331776148, 99},
		{894385949183117216, 99},
		{2775422040480279449, 99},
		{4751997750760398084, 99},
	}

	for _, c := range cases {
		heap.Insert(c.in)
		if heap.peek() != c.want {
			t.Errorf("HeapInsert(%d) == %d, wanted %d", c.in, heap.peek(), c.want)
		}
	}

	for i := 10; i >= 0; i-- {
		heap.Insert(i)
		if heap.peek() != i {
			t.Errorf("HeapInsert(%d) == %d, wanted %d", i, heap.peek(), i)
		}
	}

	for i := 0; i < 14; i++ {
		removed := heap.RemoveMin()
		expected := i
		if i > 10 {
			expected = i - 11 + 99
		}
		if removed != expected {
			fmt.Println(heap.String())
			t.Errorf("HeapRemoveMin == %d, wanted %d", removed, expected)
			t.Fail()
		}
	}
}

func TestHeapBulk(t *testing.T) {
	heap := NewHeap(func(a, b interface{}) int {
		return a.(int) - b.(int)
	})

	const count = 1000000
	arr := [count]int{}
	for i := 0; i < count; i++ {
		arr[i] = rand.Int()
		//fmt.Println(strconv.Itoa(arr[i]))
	}

	for i := 0; i < count; i++ {
		heap.Insert(arr[i])
	}

	last := heap.RemoveMin().(int)
	for i := 1; i < count; i++ {
		curr := heap.RemoveMin().(int)
		if last > curr {
			t.Errorf("HeapBulkOperation(size=%d), %d > %d", heap.size, last, curr)
		}
		last = curr
	}

	fmt.Println(heap)
}
