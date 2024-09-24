# compactor

This is the practice project to explore various golang concepts. In this project we are building Compression Tool.

NOTE: This project is part of the Coding Challenge Repo. Building A Huffman Encoder/Decoder

### Prerequisite:

Download the sample test file from [here](https://www.gutenberg.org/files/135/135-0.txt)

### Steps:

- [x] Accept filename as input and check if the file is valid and readable or not. If not throw error.

- [x] Determine the frequency of every character in the input file.

- [x] Build a binary tree from the frequency

- [x] Generate a prefix code table from the tree

- [x] Encode the text using code table.

- [x] Encode the tree - we'll need to include this in this in the output file so we can decode it.

- [x] Write the encoded tree and text to an output file.

- [x] Add the description function

- [x] Gloss-up the application using Charm tools

### Usage:

Download binary respective the device from [here](https://github.com/prashant1k99/compactor/releases)

Then extract from the archived folder.

**_Using Compactor CLI_**
You can perform 2 operations in the compactor CLI

- Compression
  - Default operation, no need to pass argument for compression apart from required flags:
    - ```~~
      ./compactor -h
      ```
- Decompression
  - For decompressing the compressed file you need to pass `dec` arg with required flags:
    - ```~~
      ./compactor dec -h
      ```

_Flags:_

- `-h`: This is the help flag to explain all the arguments and functionality of the operation.
- `-i`: [Required] This flag is required and will be pointing to the file that needs to be compressed or the compressed file which needs to be decompressed.
- `-o`: [Optional] This flag is optional, if not provided it will use the `-i` path to determine the output file
