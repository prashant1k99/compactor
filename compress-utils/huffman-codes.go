package compressutils

import (
	"errors"
	"fmt"
	"sync"
)

type HuffmanCodeTable map[rune]string

type NodePath struct {
	Node *node
	Path string
}

type HuffmanCodeChannel struct {
	Path string
	Char rune
}

func handleLeafNode(node *node, path string, huffmanCh chan<- HuffmanCodeChannel) {
	huffmanCh <- HuffmanCodeChannel{
		Char: (*node).Char(),
		Path: path,
	}
}

func handleNodes(nodeCh chan NodePath, huffmanCh chan<- HuffmanCodeChannel, wg *sync.WaitGroup) {
	defer wg.Done()

	for np := range nodeCh {
		node := np.Node
		if (*node).IsLeaf() {
			handleLeafNode(node, np.Path, huffmanCh)
		} else {
			// Handle childrens and add them to nodeCh with updatedPath
			for i, child := range (*node).Child() {
				newPath := np.Path + fmt.Sprintf("%d", i)
				nodeCh <- NodePath{
					Node: child,
					Path: newPath,
				}
			}
		}
	}
}

func TraverseBTreeToGenerateHuffmanCodes(rootNode *node, totalCodeCount int) (HuffmanCodeTable, error) {
	node := *rootNode
	if len(node.Child()) == 0 && node.IsLeaf() {
		return nil, errors.New("invalid root node: has no child and not a leaf node")
	}

	nodeCh := make(chan NodePath, 1000)
	huffmanCh := make(chan HuffmanCodeChannel)

	var wg sync.WaitGroup
	huffmanCodes := make(HuffmanCodeTable)

	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go handleNodes(nodeCh, huffmanCh, &wg)
	}

	// Start the goroutine that reads from huffmanCh
	go func() {
		processedCodes := 0
		for code := range huffmanCh {
			processedCodes++
			huffmanCodes[code.Char] = code.Path
			if processedCodes >= totalCodeCount {
				close(nodeCh)
				close(huffmanCh)
			}
		}
	}()

	// Start the root node processing
	nodeCh <- NodePath{
		Node: rootNode,
		Path: "",
	}

	// Ensure that all nodes are processed before returning
	wg.Wait()

	return huffmanCodes, nil
}
