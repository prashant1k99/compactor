package main

import (
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
	// for _, freq := range frequency {
	// 	fmt.Printf("Char: %v | count: %d \n", string(freq.Character), freq.Count)
	// }

	CreateBTreeFromFrequency(frequency)
	// for _, freq := range frequency {
	// 	fmt.Printf("Character: %v | count: %d \t", string(freq.Character), freq.Count)
	// }
	os.Exit(0)
}
