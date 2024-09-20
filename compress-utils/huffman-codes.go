package compressutils

import (
	"errors"
	"fmt"
	"sync"
)

type HuffmanCodeTable map[rune]string

type NodePath struct {
	Node *Node
	Path string
}

type HuffmanCodeChannel struct {
	Path string
	Char rune
}

func HandleLeafNode(node *Node, path string, huffmanCh chan<- HuffmanCodeChannel) {
	huffmanCh <- HuffmanCodeChannel{
		Char: (*node).Char(),
		Path: path,
	}
}

func HandleNodes(nodeCh chan NodePath, huffmanCh chan<- HuffmanCodeChannel, wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Println("Closing goroutine")

	for np := range nodeCh {
		node := np.Node
		if (*node).IsLeaf() {
			HandleLeafNode(node, np.Path, huffmanCh)
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

	nodeCh := make(chan NodePath, 1000)
	huffmanCh := make(chan HuffmanCodeChannel)

	var wg sync.WaitGroup
	huffmanCodes := make(HuffmanCodeTable)

	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go HandleNodes(nodeCh, huffmanCh, &wg)
	}
	go func() {
		fmt.Println("Scheduling cancel")
		wg.Wait()
		fmt.Println("All process finished")
		close(nodeCh)
		fmt.Println("Closed nodeCh")
		close(huffmanCh)
		fmt.Println("Closed huffmanCh")
	}()

	nodeCh <- NodePath{
		Node: rootNode,
		Path: "",
	}

	go func() {
		fmt.Println("Here:")
		for {
			code, ok := <-huffmanCh
			if !ok {
				fmt.Println("Break command called")
				break
			}
			fmt.Println("handling:", code.Char, code.Path)
			huffmanCodes[code.Char] = code.Path
		}

		fmt.Println("HERE:")
	}()

	fmt.Println("Done here:")

	return huffmanCodes, nil
}
