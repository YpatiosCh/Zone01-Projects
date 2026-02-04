package file_test

import (
	"reflect"
	"reverse/file"
	"testing"
)

func TestGetIndexes(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []int
	}{
		{
			name: "Normal case with spaces",
			input: []string{
				" ***** ",
				" ***** ",
				" ***** ",
				" ***** ",
				" ***** ",
				" ***** ",
				" ***** ",
				" ***** ",
			},
			expected: []int{0, 6},
		},
		{
			name: "Multiple line spaces",
			input: []string{
				"*** ***",
				"*** ***",
				"*** ***",
				"*** ***",
				"*** ***",
				"*** ***",
				"*** ***",
				"*** ***",
			},
			expected: []int{3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := file.GetIndexes(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestGetLetter(t *testing.T) {
	testCases := []struct {
		line     int
		expected rune
	}{
		{362, 72},  // 362 line maps to 72 which is rune of H
		{623, 101}, // 623 line maps to 101 rune of e
		{686, 108}, // 686 line maps to 108 rune of l
		{686, 108}, // 686 line maps to 108 rune of l
		{713, 111}, // 713 line maps to 111 rune of o
	}

	for _, tc := range testCases {
		t.Run(string(tc.expected), func(t *testing.T) {
			result := file.GetLetter(tc.line)
			if result != tc.expected {
				t.Errorf("For line %d, expected %c, got %c", tc.line, tc.expected, result)
			}
		})
	}
}
