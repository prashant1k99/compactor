package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	compressutils "github.com/prashant1k99/compactor/compress-utils"

	"github.com/schollz/progressbar/v3"
)

const batchSize = 1024

var (
	huffmanCodes = make(compressutils.HuffmanCodeTable)
	paddingBits  = 0
)

func writeCompressedFileMetadata(file *os.File) {
	fmt.Fprintf(file, "PaddingBits:%d\n", paddingBits)
	for key, val := range huffmanCodes {
		fmt.Fprintf(file, "%c:%s\n", key, val)
	}
	fmt.Fprintf(file, "DATA_STARTS:\n")
}

func updatePaddingBitsInMetadata(file *os.File) error {
	// Calculate the position to write
	// Assuming "PaddingBits:" is at the start of the file and we're updating the first digit after it
	position := int64(len("PaddingBits:"))

	// Seek to the position where we want to write
	_, err := file.Seek(position, 0)
	if err != nil {
		return fmt.Errorf("error seeking to position %d: %w", position, err)
	}

	// Convert PaddingBits to a string and get the first character
	paddingBitsStr := strconv.Itoa(paddingBits)
	if len(paddingBitsStr) == 0 {
		return fmt.Errorf("PaddingBits value is invalid")
	}
	byteToWrite := paddingBitsStr[0]

	// Write the single byte
	_, err = file.Write([]byte{byteToWrite})
	if err != nil {
		return fmt.Errorf("error writing byte: %w", err)
	}

	return nil
}

func convertBytesToBinary(data []byte) string {
	var binaryString string
	for _, bar := range data {
		binaryString += huffmanCodes[rune(bar)]
	}
	return binaryString
}

func convertBinaryToBytes(binaryString string, isLastBatch bool) ([]byte, string) {
	// Convert all the things to bytes if something is not wrapping it up, return it so that next batch can pick it up
	var handledBytes []byte

	// step 1 find the extra bits from binaryString and remove them from here and set them as return param
	binaryStringLen := len(binaryString)
	unprocessableBitsCount := binaryStringLen % 8
	unprocessableBits := binaryString[binaryStringLen-unprocessableBitsCount:]
	binaryString = binaryString[:binaryStringLen-unprocessableBitsCount]
	binaryStringLen = len(binaryString)

	// Step 2 Run in loop until the loop ends and make the batch of 8bits from binaryString
	for i := 0; i < binaryStringLen; i += 8 {
		chunkEndIdx := i + 8
		byteChunk := binaryString[i:chunkEndIdx]
		// Step 3: Convert string bits to processable bits
		byteVal, _ := strconv.ParseUint(byteChunk, 2, 8)
		handledBytes = append(handledBytes, byte(byteVal))
	}

	// Step4: it it's last batch, then pad the unprocessableBits to process them with additional bits and save them
	if isLastBatch && unprocessableBitsCount > 0 {
		for unprocessableBitsCount < 8 {
			paddingBits++
			unprocessableBits += "0"
			unprocessableBitsCount++
		}
		byteVal, _ := strconv.ParseUint(unprocessableBits, 2, 8)
		handledBytes = append(handledBytes, byte(byteVal))
	}

	return handledBytes, unprocessableBits
}

func CompressFile(filePath string, outputPath string) error {
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

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	readFileStat, err := file.Stat()
	if err != nil {
		fmt.Println("Error while reading input file stats:")
		return err
	}
	readFileSize := readFileStat.Size()

	bar.Describe("Generating Frequency Map")
	// First get frequency of the CompressFile
	frequncyForFile, err := compressutils.GetFrequencyForFile(filePath)
	if err != nil {
		return err
	}
	bar.Add(5)

	bar.Describe("Creating B-Tree")

	// Generate b tree and then geenrate huffman code HuffmanCodeTable
	rootNode := compressutils.CreateBTreeFromFrequency(*frequncyForFile)
	bar.Add(5)

	bar.Describe("Extracting Huffman Codes")

	totalCodeCount := len((*frequncyForFile))
	huffmanCodes, err = compressutils.TraverseBTreeToGenerateHuffmanCodes(rootNode, totalCodeCount)
	if err != nil {
		return err
	}
	bar.Add(5)

	// Open a output file for streaming
	outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	bar.Describe("Writing File Metadata")

	writeCompressedFileMetadata(outputFile)
	bar.Add(3)

	totalBytesRead := 0
	remainingBytes := ""

	bar.Describe("Compressing File")

	for {
		buffer := make([]byte, batchSize)
		byteRead, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				break
			}
		}

		if byteRead > 0 {
			// Process the batch
			totalBytesRead += byteRead
			isLastBatch := totalBytesRead == int(readFileSize)
			binaryString := convertBytesToBinary(buffer[:byteRead])
			binaryString += remainingBytes
			compressedData, remaining := convertBinaryToBytes(binaryString, isLastBatch)
			remainingBytes = remaining

			outputFile.Write(compressedData)
			// Update the CompressedPercentage variable so it can be used to show status in CLI

			progress := int(float64(totalBytesRead) / float64(readFileSize) * 82)
			bar.Set(18 + progress)
		}
	}

	// Once the padding bits is updated as per the code requirement update the metadata
	err = updatePaddingBitsInMetadata(outputFile)
	bar.Add(2)

	fmt.Printf("\nFile Compressed successfully: %s\n", outputPath)

	return err
}
