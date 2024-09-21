package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Step 1: Extract metadata from the inputFile
// Step 2: Generate ReverseHuffmanCode using metadata
// Step 3: Convert the bytes from the file to binaryString
// Step 4: Using the ReverseHuffmanCode, convert the binaryString to correct bytes
// Step 5: For last batch of the decoded, subtract the Padding Bits

type ReverseHuffmanCode map[string]rune

var (
	reverseHuffmanCode   = make(ReverseHuffmanCode)
	DecompressPercentage = 0
)

func ExtractMetadataFromFile(file *os.File) int {
	scanner := bufio.NewScanner(file)
	dataOffsetInt := 0

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimRight(line, "\r\n")
		dataOffsetInt += len(line) + 1

		if line == "DATA_STARTS:" {
			break
		}
		if strings.HasPrefix(line, "PaddingBits:") {
			paddingBits, _ = strconv.Atoi(strings.TrimPrefix(line, "PaddingBits:"))
		} else {
			parts := strings.SplitN(line, ":", 2)
			if strings.Count(line, ":") == 2 {
				parts = strings.SplitN(parts[1], ":", 2)
				reverseHuffmanCode[parts[1]] = ':'
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
				reverseHuffmanCode[parts[1]] = key
			}
		}
	}

	return dataOffsetInt
}

// func ConvertBytesToBinary()

func DecompressFile(inputFile, outputFilePath string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error while opening compressed file:")
		return err
	}
	defer file.Close()

	offsetBits := ExtractMetadataFromFile(file)

	compressedFileStats, err := file.Stat()
	if err != nil {
		fmt.Println("Error while reading file stats")
		return err
	}
	compressedFileSize := compressedFileStats.Size() - int64(offsetBits)
	fmt.Println("FileSize", compressedFileSize)
	fmt.Println("Will be outputed to:", outputFilePath)
	return nil
}
