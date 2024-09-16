package utils

import (
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestGetFrequencyCount(t *testing.T) {
	tests := []struct {
		expected Frequency
		input    string
	}{
		{
			input:    "hello",
			expected: Frequency{'h': 1, 'e': 1, 'l': 2, 'o': 1},
		},
		{
			input:    "aabbcc",
			expected: Frequency{'a': 2, 'b': 2, 'c': 2},
		},
		{
			input:    "abcABC",
			expected: Frequency{'a': 1, 'b': 1, 'c': 1, 'A': 1, 'B': 1, 'C': 1},
		},
		{
			input:    "123123",
			expected: Frequency{'1': 2, '2': 2, '3': 2},
		},
		// Edge case: Empty string
		{
			input:    "",
			expected: Frequency{},
		},
		// Edge case: String with only spaces
		{
			input:    "     ", // 5 spaces
			expected: Frequency{' ': 5},
		},
		// Edge case: String with non-alphanumeric characters
		{
			input:    "!@#$$%^&*()",
			expected: Frequency{'!': 1, '@': 1, '#': 1, '$': 2, '%': 1, '^': 1, '&': 1, '*': 1, '(': 1, ')': 1},
		},
		// Edge case: String with emojis or multi-byte characters
		{
			input:    "😀😀😁😂",
			expected: Frequency{'😀': 2, '😁': 1, '😂': 1},
		},
		// Edge case: String with a single character repeated many times
		{
			input:    "aaaaaaaaaa", // 10 'a's
			expected: Frequency{'a': 10},
		},
	}

	for _, test := range tests {
		result := GetFrequencyCount(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For input '%s', expected %v, but got %v", test.input, test.expected, result)
		}
	}
}

// Test case for sortFrequencyInAscending function
func TestSortFrequencyInAscending(t *testing.T) {
	freq := Frequency{'a': 3, 'b': 1, 'c': 2}

	expected := Frequency{'b': 1, 'c': 2, 'a': 3}

	sortedFreq := sortFrequencyInAscending(freq)
	if !reflect.DeepEqual(sortedFreq, expected) {
		t.Errorf("Expected %v, but got %v", expected, sortedFreq)
	}
}

// Test case for ProcessBatches function
func TestProcessBatches(t *testing.T) {
	// Test data
	taskCh := make(chan []byte, 1)
	freqCh := make(chan Frequency, 1)

	var wg sync.WaitGroup

	// Start the ProcessBatches goroutine
	go ProcessBatches(freqCh, taskCh, &wg)

	// Add a batch to taskCh
	taskCh <- []byte("hello")
	close(taskCh)

	// Wait for processing to complete
	wg.Wait()
	close(freqCh)

	// Expected frequency count
	expected := Frequency{'h': 1, 'e': 1, 'l': 2, 'o': 1}

	// Collect the frequency from freqCh
	result := <-freqCh

	// Compare results
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Test case for GetFrequencyForFile function
func TestGetFrequencyForFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after the test

	// Write test data to the file
	testData := "hello world"
	if _, err := tmpFile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Call GetFrequencyForFile on the temp file
	frequency, err := GetFrequencyForFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Expected frequency
	expected := Frequency{'h': 1, 'e': 1, 'l': 3, 'o': 2, ' ': 1, 'w': 1, 'r': 1, 'd': 1}

	// Check if the result matches the expected frequency
	if !reflect.DeepEqual(*frequency, expected) {
		t.Errorf("Expected %v, but got %v", expected, *frequency)
	}
}

// Test case for GetFrequencyForFile function with large file and multiple batches
func TestGetFrequencyForFile_MultipleBatches(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile_large")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after the test

	// Write a large amount of data to the file
	largeData := "aaaabbbbccccdddd" // This should span multiple batches
	if _, err := tmpFile.Write([]byte(largeData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Call GetFrequencyForFile on the temp file
	frequency, err := GetFrequencyForFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Expected frequency
	expected := Frequency{'a': 4, 'b': 4, 'c': 4, 'd': 4}

	// Check if the result matches the expected frequency
	if !reflect.DeepEqual(*frequency, expected) {
		t.Errorf("Expected %v, but got %v", expected, *frequency)
	}
}
