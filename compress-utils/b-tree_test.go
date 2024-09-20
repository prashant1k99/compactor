package compressutils

import (
	"reflect"
	"testing"
)

func TestLeafNode(t *testing.T) {
	leaf := &LeafNode{Character: 'a', Freq: 5}

	if !leaf.IsLeaf() {
		t.Error("LeafNode.IsLeaf() should return true")
	}

	if leaf.Frequency() != 5 {
		t.Errorf("LeafNode.Frequency() = %d; want 5", leaf.Frequency())
	}

	if leaf.Char() != 'a' {
		t.Errorf("LeafNode.Char() = %c; want 'a'", leaf.Char())
	}

	if leaf.Child() != nil {
		t.Error("LeafNode.Child() should return nil")
	}
}

func TestInternalNode(t *testing.T) {
	leaf1 := Node(&LeafNode{Character: 'a', Freq: 3})
	leaf2 := Node(&LeafNode{Character: 'b', Freq: 2})
	internal := &InternalNode{Children: []*Node{&leaf1, &leaf2}, Freq: 5}

	if internal.IsLeaf() {
		t.Error("InternalNode.IsLeaf() should return false")
	}

	if internal.Frequency() != 5 {
		t.Errorf("InternalNode.Frequency() = %d; want 5", internal.Frequency())
	}

	if internal.Char() != '/' {
		t.Errorf("InternalNode.Char() = %c; want 'a'", internal.Char())
	}

	if !reflect.DeepEqual(internal.Child(), []*Node{&leaf1, &leaf2}) {
		t.Error("InternalNode.Child() returned unexpected result")
	}
}

func TestPriorityQueue(t *testing.T) {
	pq := &PriorityQueue{}

	// Test empty queue
	if !pq.IsEmpty() {
		t.Error("New PriorityQueue should be empty")
	}

	if pq.Len() != 0 {
		t.Errorf("New PriorityQueue length = %d; want 0", pq.Len())
	}

	if pq.Pop() != nil {
		t.Error("Pop() on empty queue should return nil")
	}

	// Test pushing nodes (nodes should be sorted by frequency in ascending order)
	node1 := Node(&LeafNode{Character: 'a', Freq: 1})
	node2 := Node(&LeafNode{Character: 'b', Freq: 2})
	node3 := Node(&LeafNode{Character: 'c', Freq: 4})
	node4 := Node(&LeafNode{Character: 'd', Freq: 5})

	pq.Push(&node1)
	pq.Push(&node2)
	pq.Push(&node3)
	pq.Push(&node4)

	if pq.Len() != 4 {
		t.Errorf("PriorityQueue length = %d; want 4", pq.Len())
	}

	if pq.IsEmpty() {
		t.Error("PriorityQueue should not be empty after pushing")
	}

	// Test that elements are popped in ascending order of frequency
	popped := pq.Pop()
	if (*popped).Frequency() != 1 {
		t.Errorf("First Pop() frequency = %d; want 1", (*popped).Frequency())
	}

	popped = pq.Pop()
	if (*popped).Frequency() != 2 {
		t.Errorf("Second Pop() frequency = %d; want 2", (*popped).Frequency())
	}

	popped = pq.Pop()
	if (*popped).Frequency() != 4 {
		t.Errorf("Third Pop() frequency = %d; want 4", (*popped).Frequency())
	}

	popped = pq.Pop()
	if (*popped).Frequency() != 5 {
		t.Errorf("Fourth Pop() frequency = %d; want 5", (*popped).Frequency())
	}

	if !pq.IsEmpty() {
		t.Error("PriorityQueue should be empty after popping all elements")
	}
}

func TestFindCorrectInsertIndex(t *testing.T) {
	pq := &PriorityQueue{}
	node1 := Node(&LeafNode{Character: 'a', Freq: 1})
	node2 := Node(&LeafNode{Character: 'b', Freq: 3})
	node3 := Node(&LeafNode{Character: 'c', Freq: 5})

	pq.Push(&node1)
	pq.Push(&node2)
	pq.Push(&node3)

	tests := []struct {
		freq     int
		expected int
	}{
		{0, 0}, // Insert at the end
		{2, 1}, // Insert between 3 and 1
		{3, 2}, // Insert at the same frequency
		{6, 3}, // Insert at the last index
	}

	for _, tt := range tests {
		result := findCorrectInsertIndex(tt.freq, pq)

		if result != tt.expected {
			t.Errorf("findCorrectInsertIndex(%d) = %d; want %d", tt.freq, result, tt.expected)
		}
	}
}

func TestAddNewInternalLeaf(t *testing.T) {
	pq := &PriorityQueue{}
	node1 := Node(&LeafNode{Character: 'a', Freq: 1})
	node2 := Node(&LeafNode{Character: 'b', Freq: 3})
	node3 := Node(&LeafNode{Character: 'c', Freq: 5})

	pq.Push(&node1)
	pq.Push(&node2)
	pq.Push(&node3)

	newNode := Node(&LeafNode{Character: 'd', Freq: 2})
	AddNewInternalLeaf(pq, &newNode)

	expected := &PriorityQueue{}
	expected.Push(&node1)
	expected.Push(&newNode)
	expected.Push(&node2)
	expected.Push(&node3)
	if !reflect.DeepEqual(*pq, *expected) {
		t.Errorf("After AddNewInternalLeaf, pq = %v; want %v", pq, expected)
	}
}

func TestCreateBTreeFromFrequency(t *testing.T) {
	tests := []struct {
		frequency Frequency
		name      string
		wantDepth int
		wantFreq  int
	}{
		{
			name:      "Single character",
			frequency: Frequency{'a': 1},
			wantDepth: 1,
			wantFreq:  1,
		},
		{
			name:      "Two characters",
			frequency: Frequency{'a': 1, 'b': 2},
			wantDepth: 2,
			wantFreq:  3,
		},
		{
			name:      "Multiple characters",
			frequency: Frequency{'c': 1, 'b': 2, 'a': 3, 'd': 4},
			wantDepth: 4,
			wantFreq:  10,
		},
		{
			name:      "Empty frequency",
			frequency: Frequency{},
			wantDepth: 0,
			wantFreq:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := CreateBTreeFromFrequency(tt.frequency)

			if tt.wantFreq == 0 && root == nil {
				return
			}

			if tt.wantDepth == 0 {
				if root != nil {
					t.Error("Expected nil root for empty frequency")
				}
				return
			}

			if root == nil {
				t.Fatal("Unexpected nil root")
			}

			if (*root).Frequency() != tt.wantFreq {
				t.Errorf("Root frequency = %d, want %d", (*root).Frequency(), tt.wantFreq)
			}

			depth := getTreeDepth(root)
			if depth != tt.wantDepth {
				t.Errorf("Tree depth = %d, want %d", depth, tt.wantDepth)
			}

			// Validate that all characters are in the tree
			for char := range tt.frequency {
				if !containsChar(root, char) {
					t.Errorf("Tree does not contain character %c", char)
				}
			}
		})
	}
}

func getTreeDepth(node *Node) int {
	if (*node).IsLeaf() {
		return 1
	}

	maxDepth := 0
	for _, child := range (*node).Child() {
		depth := getTreeDepth(child)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth + 1
}

func containsChar(node *Node, char rune) bool {
	if (*node).IsLeaf() {
		return (*node).Char() == char
	}

	for _, child := range (*node).Child() {
		if containsChar(child, char) {
			return true
		}
	}
	return false
}
