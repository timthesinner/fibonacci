// fibonacci
package fibonacci

import (
	"bytes"
	"fmt"
	"strconv"
)

var fibonacci_numbers = [...]int{2, 3, 5, 8, 13, 21,
	34, 55, 89, 144, 233,
	377, 610, 987, 1597, 2584,
	4181, 6765, 10946, 17711, 28657,
	46368, 75025, 121393, 196418, 317811,
	514229, 832040, 1346269, 2178309, 3524578,
	5702887, 9227465, 14930352, 24157817, 39088169,
	63245986, 102334155, 165580141, 267914296, 433494437,
	701408733, 1134903170, 1836311903}

type HeapNode struct {
	value  int
	marked bool
	degree int
	parent *HeapNode
	child  *HeapNode
	left   *HeapNode
	right  *HeapNode
}

func (h *HeapNode) String() string {
	return fmt.Sprintf("Node[Degree:%d Value:%v, Marked:%t]", h.degree, h.value, h.marked)
}

type Heap struct {
	Compare       func(int, int) int
	fiboIndex     int
	fiboTarget    int
	oldFiboTarget int

	size         int
	minimum      *HeapNode
	minimumValue interface{}
}

func NewHeap(comp func(int, int) int) *Heap {
	return &Heap{Compare: comp, size: 0, fiboIndex: 0, fiboTarget: fibonacci_numbers[0], oldFiboTarget: fibonacci_numbers[0]}
}

func (h *Heap) Size() int {
	return h.size
}

func (h *Heap) String() string {
	if h.size == 0 {
		return "Heap[Size:0]"
	}

	buff := bytes.NewBufferString("Heap[Size:" + strconv.Itoa(h.size) + "\n")
	printNode(h.minimum, 1, buff)
	buff.WriteString("]")
	return buff.String()
}

func printNode(current *HeapNode, tabs int, buff *bytes.Buffer) {
	for i := 0; i < tabs; i++ {
		buff.WriteString("\t")
	}
	buff.WriteString(current.String())
	buff.WriteString("\n")
	if current.child != nil {
		printNode(current.child, tabs+1, buff)
	}

	for next := current.right; next != current; next = next.right {
		for i := 0; i < tabs; i++ {
			buff.WriteString("\t")
		}
		buff.WriteString(next.String())
		buff.WriteString("\n")
		if next.child != nil {
			printNode(next.child, tabs+1, buff)
		}
	}
}

func (h *Heap) peek() interface{} {
	if h.minimum != nil {
		return h.minimum.value
	}
	return nil
}

func (h *Heap) Insert(v int) {
	h.size++

	node := &HeapNode{value: v}
	minimum := h.minimum
	if minimum != nil {
		node.left = minimum
		node.right = minimum.right
		minimum.right = node
		node.right.left = node

		if h.Compare(node.value, minimum.value) < 0 {
			h.minimum = node
		}
	} else {
		h.minimum = node
		node.left = node
		node.right = node
	}

	if h.size%h.fiboTarget == 0 {
		h.oldFiboTarget = h.fiboTarget
		h.fiboIndex++
		h.fiboTarget = fibonacci_numbers[h.fiboIndex]
		h.consolidate()
	}
}

func swap(a, b *HeapNode) {
	v := b.value
	b.value = a.value
	a.value = v
}

func (h *Heap) consolidate() {
	if h.minimum == nil {
		return
	}

	//fmt.Println("BEFORE", h.String())

	minimum := h.minimum
	count := 1
	for current := minimum.right; current != minimum; current = current.right {
		count++
	}

	roots := make([]*HeapNode, count)
	roots[0] = minimum
	count = 1
	for current := minimum.right; current != minimum; current = current.right {
		roots[count] = current
		count++
	}

	degrees := make([]*HeapNode, 32) //TODO FIX
	for _, current := range roots {
		for degree := current.degree; degrees[degree] != nil; degree = current.degree {
			same := degrees[degree]

			//	fmt.Println("PRE MERGE", h.String(), h.minimum.String(), current.String(), same.String())
			if h.Compare(same.value, current.value) > 0 {
				merge(current, same)
				if h.minimum == same {
					h.minimum = current
				}
			} else {
				merge(same, current)
				if h.minimum == current {
					h.minimum = same
				}
				current = same
			}
			//fmt.Println(h.String())

			degrees[degree] = nil
		}
		degrees[current.degree] = current
	}

	minimum = h.minimum
	newMin := h.minimum
	for right := minimum.right; right != minimum; right = right.right {
		if h.Compare(right.value, newMin.value) < 0 {
			newMin = right
		}
	}
	h.minimum = newMin

	//fmt.Println("AFTER:", h.String())
}

func merge(parent, child *HeapNode) {
	child.right.left = child.left
	child.left.right = child.right

	child.parent = parent
	if parent.degree == 0 {
		parent.child = child
		child.left = child
		child.right = child
	} else {
		child.left = parent.child
		child.right = parent.child.right
		parent.child.right = child
		child.right.left = child
	}

	parent.degree++
}

func (h *Heap) RemoveMin() interface{} {
	oldMin := h.minimum
	if oldMin != nil {
		h.size--

		if oldMin.degree > 0 {
			child := oldMin.child
			child.parent = nil
			for right := child.right; right != child; right = right.right {
				right.parent = nil
			}

			h.minimum = child
			if oldMin.right != oldMin {
				mL := oldMin.left
				mR := oldMin.right
				cL := child.left

				mL.right = child
				child.left = mL
				cL.right = mR
				mR.left = cL
			}
			//fmt.Println("CONSOLIDATING")
			h.consolidate()
		} else if oldMin.right == oldMin {
			h.minimum = nil
		} else {
			oldMin.left.right = oldMin.right
			oldMin.right.left = oldMin.left
			h.minimum = oldMin.right
			h.consolidate()
		}

		if h.size/h.oldFiboTarget == 1 && h.size != 0 && h.oldFiboTarget == 0 {
			h.fiboIndex--
			h.fiboTarget = h.oldFiboTarget
			if h.fiboIndex > 0 {
				h.oldFiboTarget = fibonacci_numbers[h.fiboIndex]
			}
		}
		return oldMin.value
	}
	return nil
}
