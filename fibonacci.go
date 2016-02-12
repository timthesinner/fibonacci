// fibonacci
package fibonacci

import (
	"bytes"
	"fmt"
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

type Compare func(interface{}, interface{}) int

type HeapNode struct {
	value  interface{}
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
	Compare       func(interface{}, interface{}) int
	fiboIndex     int
	fiboTarget    int
	oldFiboTarget int

	size         int
	minimum      *HeapNode
	minimumValue interface{}
}

func NewHeap(comp func(interface{}, interface{}) int) *Heap {
	return &Heap{Compare: comp, fiboIndex: 0, fiboTarget: fibonacci_numbers[0], oldFiboTarget: fibonacci_numbers[0]}
}

func (h *Heap) Size() int {
	return h.size
}

func (h *Heap) String() string {
	if h.size == 0 {
		return "Heap[]"
	}

	buff := bytes.NewBufferString("Heap[\n")
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

func (h *Heap) Insert(v interface{}) {
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

	h.size++
	if h.size%h.fiboTarget == 0 {
		h.oldFiboTarget = h.fiboIndex
		h.fiboIndex += 1
		h.fiboTarget = fibonacci_numbers[h.fiboIndex]
		h.consolidate()
	}
}

func (h *Heap) consolidate() {
	numberRoots := 0
	minimum := h.minimum

	if minimum != nil {
		numberRoots = 1
		for current := minimum.right; current != minimum; current = current.right {
			numberRoots += 1
		}
	}

	nodeDegreeList := make(map[int]*HeapNode)
	current := minimum
	for numberRoots > 0 {
		degree := current.degree
		for child, ok := nodeDegreeList[degree]; ok; child, ok = nodeDegreeList[degree] {
			if child == current {
				break
			}

			parent := current
			if h.Compare(current.value, child.value) > 0 {
				temp := child
				child = parent
				parent = temp
			}

			if child == minimum {
				minimum = parent
			}
			link(child, parent)
			current = parent
			delete(nodeDegreeList, degree)
			degree++
		}

		nodeDegreeList[degree] = current
		current = current.right
		numberRoots--
	}

	newMin := minimum
	for right := minimum.right; right != minimum; right = right.right {
		if h.Compare(right.value, newMin.value) < 0 {
			newMin = right
		}
	}
	h.minimum = newMin
}

func link(child, parent *HeapNode) {
	child.left.right = child.right
	child.right.left = child.left

	child.parent = parent
	if parent.child != nil {
		child.left = parent.child
		child.right = parent.child.right
		parent.child.right = child
		child.right.left = child
	} else {
		parent.child = child
		child.left = child
		child.right = child
	}

	parent.degree++
	child.marked = false
}

func (h *Heap) RemoveMin() interface{} {
	oldMin := h.minimum
	if oldMin != nil {
		if oldMin.degree > 0 {
			child := oldMin.child
			child.parent = nil
			for right := child.right; right != child; right = right.right {
				right.parent = nil
			}

			minimum := h.minimum
			if minimum.right == minimum {
				h.minimum = child
				h.consolidate()
			} else {
				mL := minimum.left
				mR := minimum.right
				cL := child.left

				mL.right = child
				child.left = mL
				cL.right = mR
				mR.left = cL

				h.minimum = child
				h.consolidate()
			}
		} else if oldMin.right == h.minimum {
			h.minimum = nil
		} else {
			h.minimum.left.right = h.minimum.right
			h.minimum.right.left = h.minimum.left
			h.minimum = h.minimum.right
			h.consolidate()
		}

		h.size--
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
