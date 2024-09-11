package main

import (
	"bufio"
	"os"
	"sort"
	"sync"
)

type Frequency map[rune]int

func GetFrequencyFromString(data string) Frequency {
	freq := make(Frequency)
	for _, char := range data {
		freq[char]++
	}
	return freq
}

func GetFreqOfCharForBatch(data string, freqChan chan Frequency, wg *sync.WaitGroup) {
	defer wg.Done()

	freq := GetFrequencyFromString(data)
	freqChan <- freq
}

func GetFrequencyOfCharactersFromFile(f *os.File) []LeafNode {
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanRunes)

	var wg sync.WaitGroup
	freqChan := make(chan Frequency)

	const batchSize = 1024
	var batch string

	for scanner.Scan() {
		batch += scanner.Text()
		if len(batch) >= batchSize {
			wg.Add(1)
			go GetFreqOfCharForBatch(batch, freqChan, &wg)
			batch = ""
		}
	}

	if len(batch) > 0 {
		wg.Add(1)
		go GetFreqOfCharForBatch(batch, freqChan, &wg)
	}

	go func() {
		wg.Wait()
		close(freqChan)
	}()

	totalFreq := make(Frequency)
	for freq := range freqChan {
		for char, count := range freq {
			totalFreq[char] += count
		}
	}

	sortedFrequency := SortFrequencyInAscOrder(totalFreq)
	return sortedFrequency
}

func SortFrequencyInAscOrder(freq Frequency) []LeafNode {
	frequencies := make([]LeafNode, 0, len(freq))
	for char, count := range freq {
		frequencies = append(frequencies, LeafNode{Character: char, Count: count})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].Count < frequencies[j].Count // Sort in descending order
	})

	return frequencies
}
