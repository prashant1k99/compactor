package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	compressutils "github.com/prashant1k99/compactor/compress-utils"
)

const batchSize = 1024

var (
	huffmanCodes         = make(compressutils.HuffmanCodeTable)
	PaddingBits          = 0
	CompressedPercentage = 0
)

func writeCompressedFileMetadata(file *os.File) {
	fmt.Fprintf(file, "PaddingBits:%d\n", PaddingBits)
	for key, val := range huffmanCodes {
		// fmt.Printf("writing %c:%s\n", key, val)
		fmt.Fprintf(file, "%c:%s\n", key, val)
	}
	fmt.Fprintf(file, "DATA_STARTS:\n")
}

func updatePaddingBitsInMetadata(file *os.File) error {
	_, err := file.Seek(int64(len([]byte("PaddingBits:"))), 0)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(strconv.Itoa(PaddingBits)))
	return err
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
			PaddingBits++
			unprocessableBits += "0"
			unprocessableBitsCount++
		}
		fmt.Println("PaddingBits:", PaddingBits)
		byteVal, _ := strconv.ParseUint(unprocessableBits, 2, 8)
		handledBytes = append(handledBytes, byte(byteVal))
	}

	return handledBytes, unprocessableBits
}

func CompressFile(filePath string, outputPath string) error {
	// First get frequency of the CompressFile
	frequncyForFile, err := compressutils.GetFrequencyForFile(filePath)
	if err != nil {
		return err
	}
	// Generate b tree and then geenrate huffman code HuffmanCodeTable
	rootNode := compressutils.CreateBTreeFromFrequency(*frequncyForFile)
	// TODO: Need to debug the huffman code generateion, 100% ok till frequency generatiuon and logs of rootNode frequency seems to be correct
	huffmanCodes, err = compressutils.TraverseBTreeToGenerateHuffmanCodes(rootNode)
	if err != nil {
		return err
	}

	// Open a output file for streaming
	outputFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writeCompressedFileMetadata(outputFile)

	// Stream read and convert data
	// Step 1 open file for readFile
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	readFileStat, err := file.Stat()
	if err != nil {
		return err
	}
	readFileSize := readFileStat.Size()
	totalBytesRead := 0
	remainingBytes := ""

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
			CompressedPercentage = (totalBytesRead / int(readFileSize)) * 100
		}
	}

	// Once the padding bits is updated as per the code requirement update the metadata
	err = updatePaddingBitsInMetadata(outputFile)

	return err
}
