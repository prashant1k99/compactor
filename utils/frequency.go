package utils

import (
	"io"
	"os"
	"sort"
	"sync"
)

const (
	batchSize     = 1024
	maxGoroutines = 10
)

type Frequency map[rune]int

func GetFrequencyCount(data string) Frequency {
	freq := make(Frequency)
	for _, char := range data {
		freq[char]++
	}
	return freq
}

func ProcessBatches(freqCh chan Frequency, taskCh chan []byte, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for fileChunk := range taskCh {
		freq := GetFrequencyCount(string(fileChunk))
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
		go ProcessBatches(freqCh, taskCh, &wg)
	}

	for {
		buffer := make([]byte, batchSize)
		byteRead, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return nil, err
			} else {
				break
			}
		}

		if byteRead > 0 {
			taskCh <- buffer[:byteRead]
		}
	}

	close(taskCh)
	wg.Wait()

	// Collect all frequencies
	totalFreq := make(Frequency)
	for freq := range freqCh {
		for char, count := range freq {
			totalFreq[char] += count
		}
	}
	close(freqCh)
	// Sort Frequency
	totalFreq = sortFrequencyInAscending(totalFreq)

	return &totalFreq, nil
}
