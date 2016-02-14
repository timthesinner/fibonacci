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

const SKIP_CONSOLIDATION = 144

type HeapNode struct {
	value  int
	marked bool
	degree int
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
	roots        int
	minimum      *HeapNode
	minimumValue interface{}
}

func NewHeap(comp func(int, int) int) *Heap {
	return &Heap{Compare: comp, size: 0, roots: 0, fiboIndex: 0, fiboTarget: fibonacci_numbers[0], oldFiboTarget: fibonacci_numbers[0]}
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
	h.roots++

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

func (h *Heap) Consolidate() {
	h.consolidate()
}

func (h *Heap) consolidate() {
	if h.minimum == nil {
		return
	}

	roots := make([]*HeapNode, h.roots)
	current := h.minimum
	for i := 0; i < h.roots; i++ {
		roots[i] = current
		current = current.right
	}

	degrees := make([]*HeapNode, h.fiboIndex+1)
	for _, current := range roots {
		if degrees[current.degree] == nil {
			degrees[current.degree] = current
		} else {
			for degree := current.degree; degrees[degree] != nil; degree = current.degree {
				h.roots--
				same := degrees[degree]

				if h.Compare(same.value, current.value) > 0 {
					h.merge(current, same)
					if h.minimum == same {
						h.minimum = current
					}
				} else {
					h.merge(same, current)
					if h.minimum == current {
						h.minimum = same
					}
					current = same
				}

				degrees[degree] = nil
			}
			degrees[current.degree] = current
		}
	}

	h.setMinimum()
}

func (h *Heap) setMinimum() {
	minimum := h.minimum
	right := minimum.right
	for i := 0; i < h.roots; i++ {
		if h.Compare(right.value, minimum.value) < 0 {
			minimum = right
		}
		right = right.right
	}
	h.minimum = minimum
}

func (h *Heap) merge(parent, child *HeapNode) {
	child.right.left = child.left
	child.left.right = child.right

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

func (h *Heap) RemoveMin() int {
	oldMin := h.minimum
	if oldMin != nil {
		h.size--
		h.roots--

		if oldMin.degree > 0 {
			h.roots += oldMin.degree
			child := oldMin.child

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

			if h.roots <= SKIP_CONSOLIDATION {
				h.setMinimum()
			} else {
				h.consolidate()
			}
		} else if oldMin.right == oldMin {
			h.minimum = nil
		} else {
			oldMin.left.right = oldMin.right
			oldMin.right.left = oldMin.left
			h.minimum = oldMin.right
			h.consolidate()
		}

		if h.size == h.oldFiboTarget {
			h.fiboTarget = h.oldFiboTarget
			if h.fiboIndex > 0 {
				h.fiboIndex--
			}
			h.oldFiboTarget = fibonacci_numbers[h.fiboIndex]
		}

		oldMin.child = nil
		oldMin.left = nil
		oldMin.right = nil
		return oldMin.value
	}
	return 0
}
