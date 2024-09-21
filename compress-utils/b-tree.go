package compressutils

type leafMethods interface {
	IsLeaf() bool
	Char() rune
}

type internalMethods interface {
	Child() []*node
}

type node interface {
	leafMethods
	internalMethods
	Frequency() int
}

type leafNode struct {
	Character rune
	Freq      int
}

func (n *leafNode) IsLeaf() bool {
	return true
}

func (n *leafNode) Frequency() int {
	return n.Freq
}

func (n *leafNode) Char() rune {
	return n.Character
}

func (n *leafNode) Child() []*node {
	return nil
}

type internalNode struct {
	Children []*node
	Freq     int
}

func (n *internalNode) IsLeaf() bool {
	return false
}

func (n *internalNode) Frequency() int {
	return n.Freq
}

func (n *internalNode) Char() rune {
	return '/'
}

func (n *internalNode) Child() []*node {
	return n.Children
}

type PriorityQueue []*node

func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

func (pq *PriorityQueue) Push(n *node) {
	*pq = append(*pq, n)
}

func (pq *PriorityQueue) IsEmpty() bool {
	return (*pq).Len() == 0
}

func (pq *PriorityQueue) Pop() *node {
	if pq.IsEmpty() {
		return nil
	}
	item := (*pq)[0]
	*pq = (*pq)[1:]
	return item
}

func (pq *PriorityQueue) Less(i, j int) bool {
	nodeI := (*pq)[i]
	nodeJ := (*pq)[j]
	return (*nodeI).Frequency() >= (*nodeJ).Frequency()
}

func (pq *PriorityQueue) Swap(i, j int) {
	queue := (*pq)
	queue[i], queue[j] = queue[j], queue[i]
}

// Use Binary search for getting the correct insert point
func findCorrectInsertIndex(freq int, pq *PriorityQueue) int {
	low, high := 0, pq.Len()-1

	if low == high {
		return 0 // Insert at the beginning if there's only one element
	}

	for low <= high {
		mid := (low + high) / 2
		midNode := (*pq)[mid]

		if (*midNode).Frequency() < freq {
			// Move right to find the correct insertion point in ascending order
			low = mid + 1
		} else if (*midNode).Frequency() > freq {
			// Move left if freq should come before midNode
			high = mid - 1
		} else {
			// If midNode.Frequency() == freq, we want to move right to insert at the last occurrence of equal frequencies
			low = mid + 1
		}
	}

	return low // `low` now points to the correct position for insertion
}

func addNewInternalLeaf(pq *PriorityQueue, node *node) {
	indexToInsert := findCorrectInsertIndex((*node).Frequency(), pq)

	// Extend slice by 1 element
	*pq = append(*pq, nil)

	// Shift the nodes to the right
	copy((*pq)[indexToInsert+1:], (*pq)[indexToInsert:])
	// Insert the new node at correct position
	(*pq)[indexToInsert] = node
}

func CreateBTreeFromFrequency(frequency Frequency) *node {
	pq := &PriorityQueue{}

	if len(frequency) <= 0 {
		return nil
	}

	for char, freq := range frequency {
		leaf := &leafNode{
			Character: char,
			Freq:      freq,
		}
		node := node(leaf)
		pq.Push(&node)
	}

	// Create internal Nodes for the PriorityQueue
	for pq.Len() > 1 {
		minLeaf := pq.Pop()
		secondMinLeaf := pq.Pop()

		newInternalLeaf := &internalNode{
			Children: []*node{
				minLeaf,
				secondMinLeaf,
			},
			Freq: (*minLeaf).Frequency() + (*secondMinLeaf).Frequency(),
		}
		node := node(newInternalLeaf)

		addNewInternalLeaf(pq, &node)
	}
	return pq.Pop()
}
