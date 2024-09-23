package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// Step 1: Extract metadata from the inputFile
// Step 2: Generate ReverseHuffmanCode using metadata
// Step 3: Convert the bytes from the file to binaryString
// Step 4: Using the ReverseHuffmanCode, convert the binaryString to correct bytes
// Step 5: For last batch of the decoded, subtract the Padding Bits

type ReverseHuffmanCode map[string]rune

var reverseHuffmanCode = make(ReverseHuffmanCode)

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

func decompressContentInBatch(batch []byte, remainingBits string, isLastBatch bool) ([]byte, string) {
	binaryString := remainingBits + convertBytesToBinaryString(batch, 0)
	var decodedData []byte
	var currentCode string

	if isLastBatch {
		binaryString = binaryString[:len(binaryString)-(paddingBits)]
	}
	for _, bit := range binaryString {
		currentCode += string(bit)
		if char, exists := reverseHuffmanCode[currentCode]; exists {
			decodedData = append(decodedData, byte(char))
			currentCode = ""
		}
	}

	return decodedData, currentCode
}

func DecompressFile(inputFile, outputFilePath string) error {
	bar := progressbar.NewOptions(100,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("Initializing..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]█[reset]",
			SaucerHead:    "[green]█[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error while opening compressed file:")
		return err
	}
	defer file.Close()

	bar.Describe("Extracting Metadata")
	offsetBits := ExtractMetadataFromFile(file)
	bar.Add(10)

	compressedFileStats, err := file.Stat()
	if err != nil {
		fmt.Println("Error while reading file stats")
		return err
	}
	compressedFileSize := compressedFileStats.Size() - int64(offsetBits)

	// Create Output File
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Error while creating decompressedFile:")
		return err
	}
	defer outputFile.Close()

	_, err = file.Seek(int64(offsetBits), 0) // Seek to byte 100 from the beginning of the file (whence = 0)
	if err != nil {
		fmt.Println("Error seeking in file:")
		return err
	}

	bar.Describe("Decompressing File")

	buffer := make([]byte, 1024)
	remainingBits := ""
	totalBytesRead := 0
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		totalBytesRead += n
		isLastBatch := totalBytesRead == int(compressedFileSize)
		var decodedData []byte

		decodedData, remainingBits = decompressContentInBatch(buffer[:n], remainingBits, isLastBatch)

		_, err = outputFile.Write(decodedData)
		if err != nil {
			return err
		}

		if err == io.EOF {
			break
		}
		progress := int(float64(totalBytesRead) / float64(compressedFileSize) * 90)
		bar.Set(10 + progress)
	}

	fmt.Printf("\nDecompressed File Successfully: %s\n", outputFilePath)

	return nil
}
