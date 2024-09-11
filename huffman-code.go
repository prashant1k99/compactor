package main

type VisitedNode map[*Node]int

type HuffmanCodeTable struct {
	Character rune
	Code      int
}

func TraverseBTree(rootNode *Node) {
	// Create a record of visited internal nodes and path till that internal nodes
	// Always go left first, once in the root then start returning back
	// Whenever returning back to the visited node, always take the right one.
	// If returning from the right one, go 1 more node up and perform the same operation
}
