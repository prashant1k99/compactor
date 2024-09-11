# compactor

This is the practice project to explore various golang concepts. In this project we are building Compression Tool.

NOTE: This project is part of the Coding Challenge Repo. Building A Huffman Encoder/Decoder

### Prerequisite:

Download the sample test file from [here](https://www.gutenberg.org/files/135/135-0.txt)

### Steps:

- [x] Accept filename as input and check if the file is valid and readable or not. If not throw error.

- [x] Determine the frequency of every character in the input file.

- [x] Build a binary tree from the frequency

- [ ] Generate a prefix code table from the tree

- [ ] Encode the text using code table.

- [ ] Encode the tree - we'll need to include this in this in the output file so we can decode it.

- [ ] Write the encoded tree and text to an output file.
