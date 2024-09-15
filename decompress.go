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

func readCompressedFile(filePath string) ([]byte, Metadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer file.Close()

	metadata := Metadata{
		HuffmanCodeTable: make(map[rune]string),
		PaddingBits:      0,
	}

	reader := bufio.NewReader(file)
	var compressedData bytes.Buffer

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, Metadata{}, err
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "DATA_START" {
			break
		}

		if strings.HasPrefix(line, "PaddingBits:") {
			paddingBits, _ := strconv.Atoi(strings.TrimPrefix(line, "PaddingBits:"))
			metadata.PaddingBits = paddingBits
		} else {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				var key rune
				switch parts[0] {
				case "":
					key = '\n' // Empty string represents newline
				case "SPACE":
					key = ' ' // "SPACE" represents space
				default:
					key = []rune(parts[0])[0]
				}
				metadata.HuffmanCodeTable[key] = parts[1]
			}
		}

		if err == io.EOF {
			return nil, Metadata{}, fmt.Errorf("DATA_START marker not found")
		}
	}

	// Read the rest of the file as compressed data
	_, err = io.Copy(&compressedData, reader)
	if err != nil {
		return nil, Metadata{}, err
	}

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

	return binaryString
}

func decodeBinaryString(binaryString string, huffmanCodes HuffmanCodeTable) []byte {
	var decodedData []byte
	var currentCode string

	reverseHuffManCode := make(map[string]rune)
	for key, value := range huffmanCodes {
		reverseHuffManCode[value] = key
	}

	for _, bit := range binaryString {
		currentCode += string(bit)

		if char, exists := reverseHuffManCode[currentCode]; exists {
			decodedData = append(decodedData, byte(char))
			currentCode = ""
		}
	}

	return decodedData
}

func writeDecompressedContentToFile(filePath string, decompressedContent []byte) error {
	return os.WriteFile(filePath, decompressedContent, 0644)
}

func DecompressFile(compressedFilePath string, outputFilePath string) error {
	fmt.Println("DecompressingFile:", compressedFilePath)
	fmt.Println("outputFilePath", outputFilePath)

	compressedContent, metadata, err := readCompressedFile(compressedFilePath)
	if err != nil {
		return err
	}

	binaryString := convertBytesToBinaryString(compressedContent, metadata.PaddingBits)

	decompressedData := decodeBinaryString(binaryString, metadata.HuffmanCodeTable)

	err = writeDecompressedContentToFile(outputFilePath, decompressedData)
	return err
}
