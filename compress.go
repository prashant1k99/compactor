package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

type Metadata struct {
	HuffmanCodeTable
	PaddingBits int
}

func convertBytesToBinary(data []byte, huffmanCodes HuffmanCodeTable) string {
	const batchSize = 1024
	dataLen := len(data)
	numBatches := (dataLen + batchSize - 1) / batchSize // Calculate the number of batches

	results := make([]string, numBatches)
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Semaphore to limit the number of concurrent goroutines

	for i := 0; i < numBatches; i++ {
		wg.Add(1)
		sem <- struct{}{} // Acquire a semaphore slot

		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }() // Release the semaphore slot

			start := i * batchSize
			end := start + batchSize
			if end > dataLen {
				end = dataLen
			}

			var chunkBinaryString string
			for _, b := range data[start:end] {
				fmt.Printf(" Converting %c", rune(b))
				chunkBinaryString += huffmanCodes[rune(b)]
			}
			results[i] = chunkBinaryString
		}(i)
	}

	wg.Wait()

	var binaryString string
	for _, result := range results {
		binaryString += result
	}

	return binaryString
}

func padBinaryString(binaryString string) (string, int) {
	remainder := len(binaryString) % 8
	paddingBits := 0

	if remainder != 0 {
		paddingBits = 8 - remainder
		for i := 0; i < paddingBits; i++ {
			binaryString += "0"
		}
	}
	return binaryString, paddingBits
}

func convertBinaryToBytes(binaryString string) []byte {
	var byteArray []byte
	length := len(binaryString)

	for i := 0; i < length; i += 8 {
		end := i + 8
		if end > length {
			end = length
		}
		byteChunk := binaryString[i:end]

		// Pad the chunk if it's less than 8 bits
		for len(byteChunk) < 8 {
			byteChunk += "0"
		}

		byteVal, _ := strconv.ParseUint(byteChunk, 2, 8)
		byteArray = append(byteArray, byte(byteVal))
	}

	return byteArray
}

func writeCompressedFile(compressedData []byte, metadata Metadata, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("PaddingBits:%d\n", metadata.PaddingBits))
	for key, val := range metadata.HuffmanCodeTable {
		file.WriteString(fmt.Sprintf("%c:%s\n", key, val))
	}
	file.WriteString("DATA_START\n")
	savedFile, err := file.Write(compressedData)
	fmt.Println(savedFile)
	return err
}

func CompressFile(huffmanCodes HuffmanCodeTable, file *os.File) error {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	fmt.Println("here:")
	binaryString := convertBytesToBinary(fileContent, huffmanCodes)

	paddedBinaryString, paddingBits := padBinaryString(binaryString)

	compressedData := convertBinaryToBytes(paddedBinaryString)

	metadata := Metadata{
		HuffmanCodeTable: huffmanCodes,
		PaddingBits:      paddingBits,
	}

	err = writeCompressedFile(compressedData, metadata, "output.crypt")
	return err
}
