package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DecompressMetadata struct {
	HuffmanCodeTable map[string]rune
	PaddingBits      int
}

func readCompressedFile(filePath string) ([]byte, DecompressMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, DecompressMetadata{}, err
	}
	defer file.Close()

	metadata := DecompressMetadata{
		HuffmanCodeTable: make(map[string]rune),
		PaddingBits:      0,
	}

	reader := bufio.NewReader(file)
	var compressedData bytes.Buffer

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, DecompressMetadata{}, err
		}

		line = strings.TrimRight(line, "\r\n")
		fmt.Printf("Read line: %q\n", line) // Debug print

		if line == "DATA_START" {
			break
		}

		if strings.HasPrefix(line, "PaddingBits:") {
			paddingBits, _ := strconv.Atoi(strings.TrimPrefix(line, "PaddingBits:"))
			metadata.PaddingBits = paddingBits
		} else {

			parts := strings.SplitN(line, ":", 2)
			if strings.Count(line, ":") == 2 {
				parts = strings.SplitN(parts[1], ":", 2)
				metadata.HuffmanCodeTable[parts[1]] = ':'
				continue
			}
			if len(parts) == 2 {
				var key rune
				switch parts[0] {
				case "":
					key = '\n'
				case "SPACE":
					key = ' '
				default:
					key = []rune(parts[0])[0]
				}
				metadata.HuffmanCodeTable[parts[1]] = key
			}
		}

		if err == io.EOF {
			return nil, DecompressMetadata{}, fmt.Errorf("DATA_START marker not found")
		}
	}

	// Read the rest of the file as compressed data
	_, err = io.Copy(&compressedData, reader)
	if err != nil {
		return nil, DecompressMetadata{}, err
	}

	fmt.Printf("Compressed data length: %d bytes\n", compressedData.Len()) // Debug print

	return compressedData.Bytes(), metadata, nil
}

func convertBytesToBinaryString(compressedData []byte, paddingBits int) string {
	var binaryString string
	for _, b := range compressedData {
		binaryString += fmt.Sprintf("%08b", b)
	}
	if paddingBits > 0 {
		binaryString = binaryString[:len(binaryString)-paddingBits]
	}
	fmt.Printf("Binary string length: %d bits\n", len(binaryString)) // Debug print
	return binaryString
}

func decodeBinaryString(binaryString string, huffmanCodes map[string]rune) []byte {
	var decodedData []byte
	var currentCode string
	debugInterval := 100 // Adjust this value to control debug output frequency
	for i, bit := range binaryString {
		currentCode += string(bit)
		if char, exists := huffmanCodes[currentCode]; exists {
			fmt.Println("Found", char)
			decodedData = append(decodedData, byte(char))
			if i%debugInterval == 0 {
				fmt.Printf("Decoded character: %q at position %d\n", string(char), i) // Debug print
			}
			currentCode = ""
		} else if len(currentCode) > 64 { // Arbitrary limit to prevent unbounded growth
			fmt.Printf("Error: Code exceeds maximum length at position %d. Current code: %s\n", i, currentCode)
			return decodedData // Return partial data
		}

		if i%1000 == 0 {
			fmt.Printf("Processed %d bits, current decoded length: %d\n", i, len(decodedData))
		}
	}
	if currentCode != "" {
		fmt.Printf("Warning: Leftover bits at end of decoding: %s\n", currentCode)
	}
	return decodedData
}

func writeDecompressedContentToFile(filePath string, decompressedContent []byte) error {
	err := os.WriteFile(filePath, decompressedContent, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Written %d bytes to output file\n", len(decompressedContent)) // Debug print
	return nil
}

func DecompressFile(compressedFilePath string, outputFilePath string) error {
	fmt.Println("DecompressingFile:", compressedFilePath)
	fmt.Println("outputFilePath:", outputFilePath)

	compressedContent, metadata, err := readCompressedFile(compressedFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("Huffman Code Table:\n")
	for code, char := range metadata.HuffmanCodeTable {
		fmt.Printf("%q -> %q\n", code, string(char))
	}

	binaryString := convertBytesToBinaryString(compressedContent, metadata.PaddingBits)
	decompressedData := decodeBinaryString(binaryString, metadata.HuffmanCodeTable)

	fmt.Printf("Decompressed data length: %d bytes\n", len(decompressedData))
	fmt.Printf("First 100 characters of decompressed data: %q\n", string(decompressedData[:min(100, len(decompressedData))]))

	err = writeDecompressedContentToFile(outputFilePath, decompressedData)
	if err != nil {
		return err
	}

	fmt.Println("Decompression completed successfully.")
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
