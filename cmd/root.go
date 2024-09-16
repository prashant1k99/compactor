/*
Copyright © 2024 PRASHANT SINGH
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "compactor",
	Short: "Compress Single file using Huffmen Encoder",
	Run:   compressFile,
}

func compressFile(cmd *cobra.Command, args []string) {
	inputFile, err := cmd.Flags().GetString("input")
	if err != nil {
		os.Exit(1)
	}

	outputFilePath, err := cmd.Flags().GetString("output")
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Compression called", inputFile)

	fmt.Println("Output file path:", outputFilePath)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var decompressCmd = &cobra.Command{
	Use:   "dec",
	Short: "Decompress the compressed file.",
	Run:   decompressFile,
}

func decompressFile(cmd *cobra.Command, args []string) {
	fmt.Println("Decompress command called")
}

func init() {
	rootCmd.Flags().StringP("input", "i", "", "Enter the path of the file to be compressed")
	rootCmd.Flags().StringP("output", "o", "", "Enter the path for the output compressed file")
	rootCmd.MarkFlagRequired("input")
	rootCmd.AddCommand(decompressCmd)
}
