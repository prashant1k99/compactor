package utils

import (
	"fmt"
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
	wg.Add(1)
	go ProcessBatches(freqCh, taskCh, &wg)

	// Add a batch to taskCh
	taskCh <- []byte("hello")
	close(taskCh)

	// Wait for processing to complete
	wg.Wait()
	// Collect the frequency from freqCh
	result := <-freqCh
	fmt.Println(result)

	close(freqCh)

	// Expected frequency count
	expected := Frequency{'h': 1, 'e': 1, 'l': 2, 'o': 1}

	// Compare results
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetFrequencyForFile(t *testing.T) {
	// Create a temporary file for testing
	content := []byte("hello world")
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := GetFrequencyForFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("GetFrequencyForFile() error = %v", err)
	}

	expected := &Frequency{' ': 1, 'd': 1, 'e': 1, 'h': 1, 'l': 3, 'o': 2, 'r': 1, 'w': 1}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetFrequencyForFile() = %v, want %v", result, expected)
	}
}
