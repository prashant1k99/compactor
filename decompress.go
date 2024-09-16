package main

import (
	"bufio"
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

func readCompressedFileMatadata(file *os.File) (DecompressMetadata, int64, error) {
	metadata := DecompressMetadata{
		HuffmanCodeTable: make(map[string]rune),
	}

	reader := bufio.NewReader(file)
	dataOffsetInt := int64(0)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return DecompressMetadata{}, 0, err
		}

		line = strings.TrimRight(line, "\r\n")
		fmt.Printf("Read line: %q\n", line)

		dataOffsetInt += int64(len(line) + 1)

		if line == "DATA_START" {
			// dataOffsetInt, _ = file.Seek(0, io.SeekCurrent)
			// fmt.Println(dataOffsetInt)
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
			return DecompressMetadata{}, 0, fmt.Errorf("DATA_START marker not found")
		}
	}

	return metadata, dataOffsetInt, nil
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

func decompressContentInBatch(batch []byte, remainingBits string, metadata *DecompressMetadata, isLastBatch bool) ([]byte, string) {
	binaryString := remainingBits + convertBytesToBinaryString(batch, 0)
	fmt.Println(binaryString)
	var decodedData []byte
	var currentCode string

	fmt.Println("TotalLength", len(batch))
	if isLastBatch {
		binaryString = binaryString[:len(binaryString)-(*&metadata.PaddingBits)]
	}
	for _, bit := range binaryString {
		currentCode += string(bit)
		if char, exists := (*metadata).HuffmanCodeTable[currentCode]; exists {
			decodedData = append(decodedData, byte(char))
			currentCode = ""
		}
	}

	return decodedData, currentCode
}

func DecompressFile(compressedFilePath string, outputFilePath string) error {
	compressedFile, err := os.Open(compressedFilePath)
	if err != nil {
		return err
	}
	defer compressedFile.Close()

	metadata, offsetInt, err := readCompressedFileMatadata(compressedFile)
	if err != nil {
		return err
	}

	compressedFileStats, err := compressedFile.Stat()
	if err != nil {
		return err
	}
	compressedFileTotalSize := compressedFileStats.Size()

	_, err = compressedFile.Seek(offsetInt, io.SeekStart)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	remainingBits := ""

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	totalBytesRead := int64(0)

	for {
		n, err := compressedFile.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		totalBytesRead += int64(n)
		isLastBatch := totalBytesRead == compressedFileTotalSize-offsetInt
		var decodedData []byte

		decodedData, remainingBits = decompressContentInBatch(buffer[:n], remainingBits, &metadata, isLastBatch)

		_, err = outputFile.Write(decodedData)
		if err != nil {
			return err
		}

		if len(remainingBits) > 0 && len(remainingBits) != metadata.PaddingBits {
			fmt.Printf("Warning: %d bits remained unprocessed at the end\n", len(remainingBits))
		}

		if err == io.EOF {
			break
		}
	}
	fmt.Println("Decompression completed successfully.")
	return nil
}
