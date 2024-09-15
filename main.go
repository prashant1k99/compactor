package main

import (
	"flag"
	"fmt"
	"os"
)

func ReadFile(filename string) (*os.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func main() {
	decompres := flag.Bool("dec", false, "Decompress execution")
	flag.Parse()
	if *decompres {
		err := DecompressFile("output.crypt", "test1.txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Successfully decompressed")
		os.Exit(0)
	}
	fileName := "smallString.txt"
	file, err := ReadFile(fileName)
	if err != nil {
		fmt.Printf("unable to read file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	frequency := GetFrequencyOfCharactersFromFile(file)
	for _, freq := range frequency {
		fmt.Printf("Char: %v | count: %d \n", string(freq.Character), freq.Count)
	}

	rootNode := CreateBTreeFromFrequency(frequency)

	huffmanCodes, err := TraverseBTree(&rootNode)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	// for _, freq := range frequency {
	// 	fmt.Printf("Character: %v | count: %d \t", string(freq.Character), freq.Count)
	// }

	// Create the metadata for the file and encrypt the contents of the fileS
	f, err := ReadFile(fileName)
	if err != nil {
		fmt.Println("unable to read file", err)
		os.Exit(1)
	}
	err = CompressFile(huffmanCodes, f)
	if err != nil {
		fmt.Println("failed to compress file:", err)
		os.Exit(1)
	}
	// Decrypt the encrypted file
	os.Exit(0)
}
