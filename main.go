package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type Frequency map[rune]int

func ReadFile(filename string) (*os.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GetFrequencyOfCharactersForBatch(data string, freqChan chan Frequency, wg *sync.WaitGroup) {
	defer wg.Done()

	freq := make(Frequency)
	for _, char := range data {
		freq[char]++
	}
	freqChan <- freq
}

func GetFrequencyOfCharactersFromFile(f *os.File) Frequency {
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
			go GetFrequencyOfCharactersForBatch(batch, freqChan, &wg)
			batch = ""
		}
	}

	if len(batch) > 0 {
		wg.Add(1)
		go GetFrequencyOfCharactersForBatch(batch, freqChan, &wg)
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

	return totalFreq
}

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("invalid input: filename to encode missing")
		os.Exit(1)
	}

	fmt.Println(args[0])
	file, err := ReadFile(args[0])
	if err != nil {
		fmt.Printf("unable to read file: %v", err)
		os.Exit(1)
	}
	frequency := GetFrequencyOfCharactersFromFile(file)
	for char, count := range frequency {
		fmt.Printf("Character %v | count %d \t", string(char), count)
	}
	os.Exit(0)
}
