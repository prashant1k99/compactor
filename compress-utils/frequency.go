package compressutils

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
)

const (
	batchSize     = 1024
	maxGoroutines = 10
)

var FrequencyProgress = 0

type Frequency map[rune]int

func getFrequencyCount(data string) Frequency {
	freq := make(Frequency)
	for _, char := range data {
		freq[char]++
	}
	return freq
}

func processBatches(freqCh chan Frequency, taskCh chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	for fileChunk := range taskCh {
		freq := getFrequencyCount(string(fileChunk))
		freqCh <- freq
	}
}

func sortFrequencyInAscending(freq Frequency) Frequency {
	type RuneFreq struct {
		Key   rune
		Value int
	}
	// Step 1: Convert the Frequency map to a slice of RuneFreq pairs
	var freqSlice []RuneFreq
	for key, value := range freq {
		freqSlice = append(freqSlice, RuneFreq{Key: key, Value: value})
	}

	// Step 2: Sort the slice based on the frequency values
	sort.Slice(freqSlice, func(i, j int) bool {
		return freqSlice[i].Value < freqSlice[j].Value
	})

	// Step 3: Convert the sorted slice back to a Frequency map
	sortedFreq := make(Frequency)
	for _, kv := range freqSlice {
		sortedFreq[kv.Key] = kv.Value
	}

	return sortedFreq
}

func GetFrequencyForFile(filePath string) (*Frequency, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wg sync.WaitGroup
	freqCh := make(chan Frequency)
	taskCh := make(chan []byte, maxGoroutines)

	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go processBatches(freqCh, taskCh, &wg)
	}

	go func() {
		wg.Wait()
		close(freqCh)
	}()

	go func() {
		defer close(taskCh)

		for {
			buffer := make([]byte, batchSize)
			byteRead, err := file.Read(buffer)
			fmt.Println("Bytes Read:", byteRead)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading file:", err)
					return
				} else {
					break
				}
			}

			if byteRead > 0 {
				taskCh <- buffer[:byteRead]
			}
		}
	}()
	// Collect all frequencies
	totalFreq := make(Frequency)
	for freq := range freqCh {
		for char, count := range freq {
			totalFreq[char] += count
		}
	}

	// Sort Frequency
	totalFreq = sortFrequencyInAscending(totalFreq)

	return &totalFreq, nil
}
