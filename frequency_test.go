package main

import (
	"testing"
)

func TestFrequencyGenerator(t *testing.T) {
	tests := []struct {
		input    string
		expected []LeafNode
	}{
		{
			"abbbbbbcc",
			[]LeafNode{
				{'a', 1},
				{'c', 2},
				{'b', 6},
			},
		},
		{
			"aaaaaaaaa",
			[]LeafNode{
				{'a', 9},
			},
		},
	}

	for _, test := range tests {
		freq := GetFrequencyFromString(test.input)
		result := SortFrequencyInAscOrder(freq)

		if !equal(result, test.expected) {
			t.Errorf("For input '%s', expected %v, but got %v", tests[0].input, tests[0].expected, result)
		}
	}
}

func equal(a, b []LeafNode) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
