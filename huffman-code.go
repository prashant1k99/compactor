package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type HuffmanCodeTable map[rune]string

var (
	HuffmanCodes = make(HuffmanCodeTable)
	mu           sync.Mutex
)

type NodePath struct {
	Node *Node
	Path string
}

func HandleLeafNode(node *Node, path string) {
	fmt.Println("Processing", (*node).Char())
	mu.Lock()
	defer mu.Unlock()
	HuffmanCodes[(*node).Char()] = path
	fmt.Println("Processed", (*node).Char())
}

func TraverseBTree(rootNode *Node) (HuffmanCodeTable, error) {
	if len((*rootNode).Child()) == 0 && (*rootNode).IsLeaf() {
		return HuffmanCodeTable{}, fmt.Errorf("no children node found on the root node and is not leaf")
	}

	nodeChan := make(chan NodePath, 100)
	done := make(chan struct{})
	var wg sync.WaitGroup
	nodeCount := int32(0)
	processedCount := int32(0)

	// Start worker goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for np := range nodeChan {
				fmt.Printf("Worker %d processing node: %v\n", id, (*np.Node).Frequency())
				HandleNode(np, nodeChan, &wg, &nodeCount, &processedCount)
			}
			fmt.Printf("Worker %d exiting\n", id)
		}(i)
	}

	// Start a goroutine to close the done channel when all workers are finished
	go func() {
		wg.Wait()
		fmt.Println("All workers finished, closing done channel")
		close(done)
	}()

	// Send the root node to start the process
	fmt.Println("Sending root node to channel")
	wg.Add(1)
	atomic.AddInt32(&nodeCount, 1)
	nodeChan <- NodePath{
		Node: rootNode,
		Path: "",
	}

	// Wait for either all nodes to be processed or a timeout
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			fmt.Println("All processing completed")
			close(nodeChan)
			return HuffmanCodes, nil
		case <-timeout:
			fmt.Printf("Traversal timed out. Processed %d out of %d nodes\n", atomic.LoadInt32(&processedCount), atomic.LoadInt32(&nodeCount))
			close(nodeChan)
			return nil, fmt.Errorf("traversal timed out")
		case <-ticker.C:
			processed := atomic.LoadInt32(&processedCount)
			total := atomic.LoadInt32(&nodeCount)
			fmt.Printf("Progress: Processed %d out of %d nodes\n", processed, total)
			if processed == total && total > 0 {
				fmt.Println("All nodes processed, closing channels")
				close(nodeChan)
				<-done // Wait for workers to finish
				return HuffmanCodes, nil
			}
		}
	}
}

func HandleNode(np NodePath, nodeChan chan<- NodePath, wg *sync.WaitGroup, nodeCount *int32, processedCount *int32) {
	defer wg.Done()
	atomic.AddInt32(processedCount, 1)
	node := np.Node

	if (*node).IsLeaf() {
		fmt.Printf("Handling leaf node: %v\n", (*node).Char())
		HandleLeafNode(node, np.Path)
		return
	}

	childNodes := (*node).Child()
	fmt.Printf("Processing node with %d children, frequency: %d\n", len(childNodes), (*node).Frequency())

	if len(childNodes) == 0 {
		fmt.Printf("Warning: non-leaf node with no children found: %v\n", (*node).Frequency())
		return
	}

	for i, child := range childNodes {
		wg.Add(1)
		atomic.AddInt32(nodeCount, 1)
		newPath := np.Path + fmt.Sprintf("%d", i)
		fmt.Printf("Sending child node to channel: frequency %v with path %s\n", child.Frequency(), newPath)
		nodeChan <- NodePath{Node: &child, Path: newPath}
	}
}
