package compressutils

import (
	"reflect"
	"sync"
	"testing"
)

func TestHandleLeafNode(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		node     node
		expected HuffmanCodeChannel
	}{
		{
			name: "Single leaf node",
			node: &leafNode{Character: 'a', Freq: 1},
			path: "010100",
			expected: HuffmanCodeChannel{
				Char: 'a',
				Path: "010100",
			},
		},
		{
			name:     "Multiple leaf nodes",
			node:     &leafNode{Character: 'b', Freq: 2},
			path:     "10",
			expected: HuffmanCodeChannel{Char: 'b', Path: "10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			huffmanCh := make(chan HuffmanCodeChannel, 2)

			handleLeafNode(&tt.node, tt.path, huffmanCh)

			huffmanCodes := <-huffmanCh

			if !reflect.DeepEqual(huffmanCodes, tt.expected) {
				t.Errorf("HandleLeafNode() = %v, want %v", huffmanCodes, tt.expected)
			}
		})
	}
}

func TestHandleNodes(t *testing.T) {
	nodeA := node(&leafNode{
		Character: 'a',
		Freq:      1,
	})
	nodeB := node(&leafNode{
		Character: 'b',
		Freq:      2,
	})
	nodeC := node(&leafNode{
		Character: 'c',
		Freq:      '1',
	})

	twoLeafNode := node(&internalNode{
		Children: []*node{
			&nodeB,
			&nodeC,
		},
		Freq: 3,
	})
	tests := []struct {
		expected HuffmanCodeTable
		name     string
		nodes    []NodePath
	}{
		{
			name: "Single leaf node",
			nodes: []NodePath{
				{Node: &nodeA, Path: "0101"},
			},
			expected: HuffmanCodeTable{'a': "0101"},
		},
		{
			name: "Internal node with two leaf children",
			nodes: []NodePath{
				{Node: &twoLeafNode, Path: ""},
			},
			expected: HuffmanCodeTable{'b': "0", 'c': "1"},
		},
	}

	for _, tt := range tests {
		nodeCh := make(chan NodePath, 3)
		huffmanCh := make(chan HuffmanCodeChannel)
		var wg sync.WaitGroup

		wg.Add(1)
		go handleNodes(nodeCh, huffmanCh, &wg)

		nodeCh <- tt.nodes[0]

		go func() {
			wg.Wait()
			close(nodeCh)
		}()

		huffmanCodes := make(HuffmanCodeTable)
		go func() {
			for code := range huffmanCh {
				huffmanCodes[code.Char] = code.Path
			}
			close(huffmanCh)
			t.Run(tt.name, func(t *testing.T) {
				if !reflect.DeepEqual(huffmanCodes, tt.expected) {
					t.Errorf("HandleNodes() = %v, want %v", huffmanCodes, tt.expected)
				}
			})
		}()
	}
}

func TestTraverseBTree(t *testing.T) {
	nodeA := node(&leafNode{
		Character: 'a',
		Freq:      2,
	})
	nodeB := node(&leafNode{
		Character: 'b',
		Freq:      1,
	})

	nodeC := node(&leafNode{
		Character: 'c',
		Freq:      5,
	})

	twoLeafNode := node(&internalNode{
		Children: []*node{
			&nodeA,
			&nodeB,
		},
		Freq: 3,
	})

	tests := []struct {
		rootNode           node
		expected           HuffmanCodeTable
		name               string
		codeToBeCalculated int
		expectedErr        bool
	}{
		{
			name:               "Valid tree",
			rootNode:           twoLeafNode,
			expected:           HuffmanCodeTable{'a': "0", 'b': "1"},
			codeToBeCalculated: 2,
			expectedErr:        false,
		},
		{
			name:               "Invalid root node (leaf node)",
			rootNode:           &leafNode{Character: 'a', Freq: 1},
			expected:           nil,
			codeToBeCalculated: 0,
			expectedErr:        true,
		},
		{
			name: "Complex tree",
			rootNode: &internalNode{
				Children: []*node{
					&twoLeafNode,
					&nodeC,
				},
				Freq: 9,
			},
			expected:           HuffmanCodeTable{'a': "00", 'b': "01", 'c': "1"},
			codeToBeCalculated: 3,
			expectedErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TraverseBTreeToGenerateHuffmanCodes(&tt.rootNode, tt.codeToBeCalculated)
			if (err != nil) != tt.expectedErr {
				t.Errorf("TraverseBTree() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TraverseBTree() = %v, want %v", result, tt.expected)
			}
		})
	}
}
