package compressutils

import (
	"errors"
	"fmt"
	"sync"
)

type HuffmanCodeTable map[rune]string

var (
	HuffmanCodes = make(HuffmanCodeTable)
	mu           sync.RWMutex
)

type NodePath struct {
	Node *Node
	Path string
}

func HandleLeafNode(node *Node, path string) {
	mu.Lock()
	defer mu.Unlock()
	HuffmanCodes[(*node).Char()] = path
}

func HandleNodes(nodeCh chan NodePath, wg *sync.WaitGroup) {
	defer wg.Done()

	for np := range nodeCh {
		node := np.Node
		if (*node).IsLeaf() {
			HandleLeafNode(node, np.Path)
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

func TraverseBTreeToGenerateHuffmanCodes(rootNode *Node) (HuffmanCodeTable, error) {
	node := *rootNode
	if len(node.Child()) == 0 && node.IsLeaf() {
		return nil, errors.New("invalid root node: has no child and not a leaf node")
	}

	nodeCh := make(chan NodePath)

	var wg sync.WaitGroup

	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go HandleNodes(nodeCh, &wg)
	}

	go func() {
		wg.Wait()
		close(nodeCh)
	}()

	nodeCh <- NodePath{
		Node: rootNode,
		Path: "",
	}

	return HuffmanCodes, nil
}
