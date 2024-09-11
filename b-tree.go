package main

import "fmt"

type Node interface {
	Frequency() int
	IsLeaf() bool
}

type LeafNode struct {
	Character rune
	Count     int
}

func (n LeafNode) Frequency() int {
	return n.Count
}

func (n LeafNode) IsLeaf() bool {
	return true
}

type InternalNode struct {
	Children []Node
	Count    int
}

func (n InternalNode) Frequency() int {
	return n.Count
}

func (n InternalNode) IsLeaf() bool {
	return false
}

type PriorityQueue []Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(Node)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() Node {
	old := *pq
	item := old[0]
	*pq = old[1:]
	return item
}

func (pq *PriorityQueue) Peek() Node {
	return (*pq)[0]
}

func (pq *PriorityQueue) IsEmpty() bool {
	return pq.Len() <= 1
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Frequency() > pq[j].Frequency()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func PlaceNewInternalLeafInPlace(pq *PriorityQueue, node *InternalNode) {
	fmt.Println("processing", node.Frequency())
	pq.Push(node)
	i := pq.Len() - 1
	for i > 0 {
		if pq.Less(i-1, i) {
			pq.Swap(i, i-1)
		}
		if (*pq)[i].Frequency() <= node.Frequency() {
			break
		}
		i--
	}
}

func CreateBTreeFromFrequency(freq []LeafNode) {
	// We need to save all the
	pq := &PriorityQueue{}

	for _, leaf := range freq {
		pq.Push(leaf)
	}

	for !pq.IsEmpty() {
		minLeafNode := pq.Pop()
		secondMinLeafNode := pq.Pop()
		fmt.Println("min:", minLeafNode.Frequency())
		fmt.Println("secondMin:", secondMinLeafNode.Frequency())
		newInternalNode := InternalNode{
			Children: []Node{
				minLeafNode,
				secondMinLeafNode,
			},
			Count: minLeafNode.Frequency() + secondMinLeafNode.Frequency(),
		}
		PlaceNewInternalLeafInPlace(pq, &newInternalNode)
	}
	fmt.Println("root node", (*pq)[0])
}
