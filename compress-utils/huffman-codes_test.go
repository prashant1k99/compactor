package compressutils

import (
	"reflect"
	// "sync"
	"testing"
)

func TestHandleLeafNode(t *testing.T) {
	tests := []struct {
		expected HuffmanCodeTable
		name     string
		node     Node
		path     string
	}{
		{
			name:     "Single leaf node",
			node:     &LeafNode{Character: 'a', Freq: 1},
			path:     "0",
			expected: HuffmanCodeTable{'a': "0"},
		},
		{
			name:     "Multiple leaf nodes",
			node:     &LeafNode{Character: 'b', Freq: 2},
			path:     "10",
			expected: HuffmanCodeTable{'b': "10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset HuffmanCodes before each sub-test
			mu.Lock()
			HuffmanCodes = make(HuffmanCodeTable)
			mu.Unlock()

			HandleLeafNode(&tt.node, tt.path)

			mu.RLock()
			if !reflect.DeepEqual(HuffmanCodes, tt.expected) {
				t.Errorf("HandleLeafNode() = %v, want %v", HuffmanCodes, tt.expected)
			}
			mu.RUnlock()
		})
	}
}

// func TestHandleNodes(t *testing.T) {
// 	// Reset HuffmanCodes before each test
//
// 	nodeA := Node(&LeafNode{
// 		Character: 'a',
// 		Freq:      1,
// 	})
// 	nodeB := Node(&LeafNode{
// 		Character: 'b',
// 		Freq:      2,
// 	})
// 	nodeC := Node(&LeafNode{
// 		Character: 'c',
// 		Freq:      '1',
// 	})
//
// 	twoLeafNode := Node(&InternalNode{
// 		Children: []*Node{
// 			&nodeB,
// 			&nodeC,
// 		},
// 		Freq: 3,
// 	})
// 	tests := []struct {
// 		expected HuffmanCodeTable
// 		name     string
// 		nodes    []NodePath
// 	}{
// 		{
// 			name: "Single leaf node",
// 			nodes: []NodePath{
// 			{Node: &nodeA, Path: "0101"},
// 		},
// 		expected: HuffmanCodeTable{'a': "0101"},
// 	},
// 	{
// 		name: "Internal node with two leaf children",
// 		nodes: []NodePath{
// 			{Node: &twoLeafNode, Path: ""},
// 		},
// 		expected: HuffmanCodeTable{'b': "0", 'c': "1"},
// 	},
// }
//
// for _, tt := range tests {
// 	t.Run(tt.name, func(t *testing.T) {
// 		// mu.Lock()
// 		HuffmanCodes = make(HuffmanCodeTable)
// 		// mu.Unlock()
//
// 		nodeCh := make(chan NodePath, len(tt.nodes))
// 		var wg sync.WaitGroup
// 		wg.Add(1)
//
// 		go HandleNodes(nodeCh, &wg)
//
// 		for _, np := range tt.nodes {
// 			nodeCh <- np
// 		}
// 		go func() {
// 			wg.Wait()
// 			close(nodeCh)
// 			}()
//
// 			if !reflect.DeepEqual(HuffmanCodes, tt.expected) {
// 				t.Errorf("HandleNodes() = %v, want %v", HuffmanCodes, tt.expected)
// 			}
// 		})
// 	}
// }

//
// func TestTraverseBTree(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		rootNode    Node
// 		expected    HuffmanCodeTable
// 		expectedErr bool
// 	}{
// 		{
// 			name: "Valid tree",
// 			rootNode: &InternalNode{
// 				Children: []*Node{
// 					{&LeafNode{Character: 'a', Freq: 2}},
// 					{&LeafNode{Character: 'b', Freq: 1}},
// 				},
// 				Freq: 3,
// 			},
// 			expected:    HuffmanCodeTable{'a': "0", 'b': "1"},
// 			expectedErr: false,
// 		},
// 		{
// 			name:        "Invalid root node (leaf node)",
// 			rootNode:    &LeafNode{Character: 'a', Freq: 1},
// 			expected:    nil,
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Complex tree",
// 			rootNode: &InternalNode{
// 				Children: []*Node{
// 					{&InternalNode{
// 						Children: []*Node{
// 							{&LeafNode{Character: 'a', Freq: 3}},
// 							{&LeafNode{Character: 'b', Freq: 2}},
// 						},
// 						Freq: 5,
// 					}},
// 					{&LeafNode{Character: 'c', Freq: 4}},
// 				},
// 				Freq: 9,
// 			},
// 			expected:    HuffmanCodeTable{'a': "00", 'b': "01", 'c': "1"},
// 			expectedErr: false,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Reset HuffmanCodes before each test
// 			HuffmanCodes = make(HuffmanCodeTable)
//
// 			result, err := TraverseBTree(&tt.rootNode)
//
// 			if (err != nil) != tt.expectedErr {
// 				t.Errorf("TraverseBTree() error = %v, expectedErr %v", err, tt.expectedErr)
// 				return
// 			}
//
// 			if !reflect.DeepEqual(result, tt.expected) {
// 				t.Errorf("TraverseBTree() = %v, want %v", result, tt.expected)
// 			}
// 		})
// 	}
// }
