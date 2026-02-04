package asciiart

import (
	//"fmt"
	"os"
	"strings"
	"runtime"
)

// LocateLines locates the lines to start printing each word seperately.
// Returns a two-dimensional slice of integers where each index represents a word and
// its inner indexes represent the lines of each letter of the word
// For example: "Hello\nworld" will become -> [[362 623 686 686 713] [785 713 740 686 614]]
// Note that we seperate words only when we encounter "\n" and not " "

func LocateLines(input string) [][]int {
	// Split into different words when encountering "\n"
	array := strings.Split(input, "\\n")

	// Initialize an empty slice of slices to store the lines we need to start from
	lines := make([][]int, len(array))

	for i, word := range array {
		// Initialize a new slice for the current word
		lines[i] = []int{}
		for _, letter := range word {
			// Convert the letter to its integer representation (Unicode code point)
			line := int(letter)

			// Apply the mathematical formula to find the line number
			lineFromText := (line-32)*9 + 2

			// Append the result to the slice for the current word
			lines[i] = append(lines[i], lineFromText)
		}
	}
	//fmt.Println(lines)

	return lines
}

// PrintLinesFromArray reads the specified lines from a text file and prints them
// Before continuing to the next index of the two-dimensional slice, it will print the next 8 lines of each index in the inner array
func PrintLinesFromArray(filename string, lineNumbers [][]int) ([]string, error) {
	// Read the entire file into memory
	content, err := os.ReadFile(filename)

	if err != nil {
		return []string{}, err
	}
	var lines []string
	if runtime.GOOS == "windows" {
        // If running on Windows, normalize \r\n to \n
        normalizedContent := strings.ReplaceAll(string(content), "\r\n", "\n")
        lines = strings.Split(normalizedContent, "\n")
    } else {
        // For Unix-like systems (Linux, macOS), just split using \n
        lines = strings.Split(string(content), "\n")
    }

	var outputLines []string

	// // Iterate over each array of line numbers
	for _, lineArray := range lineNumbers {
		for i := 0; i < 8; i++ {
			// if we find an empty array, we break the loop after printing new line
			if len(lineArray) == 0 {
				outputLines = append(outputLines, "")
				break
			}
			var currentLine string
			// Print the next 8 lines of each index in the array
			for j := 0; j < len(lineArray); j++ {
				currentLine += lines[lineArray[j]-1]
				lineArray[j]++
			}
			outputLines = append(outputLines, currentLine)
		}
	}
	return outputLines, nil
}
