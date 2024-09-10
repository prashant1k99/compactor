package main

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
	Children []*Node
	Count    int
}

func (n InternalNode) Frequency() int {
	return n.Count
}

func (n InternalNode) IsLeaf() bool {
	return false
}

type PriorityQueue []*Node

func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

func (pq *PriorityQueue) Push(node *Node) {
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() Node {
	item := (*pq)[0]
	*pq = (*pq)[1:]
	return *item
}

func (pq *PriorityQueue) Peek() *Node {
	return (*pq)[0]
}

func (pq *PriorityQueue) IsEmpty() bool {
	return pq.Len() <= 0
}

func (pq *PriorityQueue) Less(i, j int) bool {
	// Compare the frequency of the two nodes
	return (*(*pq)[i]).Frequency() > (*(*pq)[j]).Frequency()
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func CreateBTreeFromFrequency(freq []LeafNode) {
}
