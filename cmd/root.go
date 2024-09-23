/*
Copyright © 2024 PRASHANT SINGH
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "compactor",
	Short: "Compress Single file using Huffmen Encoder",
	Run:   compressFile,
}

var rootCmdHelpTemplate = `{{with .Short}}{{. | trimTrailingWhitespaces}}{{end}}

Usage:
  {{.UseLine}}

  # Decompress a file
  compactor dec


Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Description:
  This command compresses a single file using Huffman encoding. You need to provide the input file path and optionally the output file path.
  Default output path is whatever the folder path for input file

Examples:
  # Compress a file
  compactor -i input.txt -o output.crypt

  # Compress a file with default output path
  compactor -i input.txt

`

// Custom help template for decompressCmd
var decompressCmdHelpTemplate = `{{with .Short}}{{. | trimTrailingWhitespaces}}{{end}}

Usage:
  {{.UseLine}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Description:
  This command decompresses a file that was compressed using Huffman encoding. You need to provide the input compressed file path and optionally the output file path.
  Default output path is whatever the folder path for input file

Examples:
  # Decompress a file
  compactor dec -i input.crypt -o output.txt

  # Decompress a file with default output path
  compactor dec -i input.crypt

`

func compressFile(cmd *cobra.Command, args []string) {
	inputFile, err := cmd.Flags().GetString("input")
	if err != nil {
		os.Exit(1)
	}

	outputFilePath, err := cmd.Flags().GetString("output")
	if err != nil {
		os.Exit(1)
	}

	if outputFilePath == "" {
		inputDir := filepath.Dir(inputFile)
		outputFilePath = inputDir
	}
	inputFileName := filepath.Base(inputFile)
	outputFilePath = filepath.Join(outputFilePath, inputFileName+".crypt")

	err = CompressFile(inputFile, outputFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
	inputFile, err := cmd.Flags().GetString("input")
	if err != nil {
		os.Exit(1)
	}

	outputFilePath, err := cmd.Flags().GetString("output")
	if err != nil {
		os.Exit(1)
	}

	if outputFilePath == "" {
		inputDir := filepath.Dir(inputFile)
		outputFilePath = inputDir
	}
	inputFileName := filepath.Base(inputFile)
	inputFileNameWithoutExt := strings.TrimSuffix(inputFileName, filepath.Ext(inputFileName))
	outputFilePath = filepath.Join(outputFilePath, inputFileNameWithoutExt)

	err = DecompressFile(inputFile, outputFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("input", "i", "", "Enter the path of the file to be compressed")
	rootCmd.Flags().StringP("output", "o", "", "Enter the path for the output compressed file")
	rootCmd.Flags().BoolP("help", "h", false, "Show help for all the options")
	rootCmd.MarkFlagRequired("input")

	decompressCmd.Flags().StringP("input", "i", "", "Enter file path of Compressed file")
	decompressCmd.Flags().StringP("output", "o", "", "Enter path for decompressed file")
	decompressCmd.MarkFlagRequired("input")
	decompressCmd.Flags().BoolP("help", "h", false, "Show help for all the options")

	rootCmd.SetHelpTemplate(rootCmdHelpTemplate)
	decompressCmd.SetHelpTemplate(decompressCmdHelpTemplate)

	rootCmd.AddCommand(decompressCmd)
}
