package tools

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// ToFileOrNotToFile checks whether a "--output" flag is being used to determine
// if the output is going to print to terminal or write to file.
// It will return an error if the user incorrectly tried to type the flag
func ToFileOrNotToFile(s string) bool {
	if strings.HasPrefix(s, "--output=") && strings.HasSuffix(s, ".txt") {
		return true
	}
	return false
}

func PrintUsage() {
	fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
	fmt.Println("\nEX: go run . --output=<fileName.txt> something standard")
}

// LocateLines locates the lines to start printing each word seperately.
// Returns a two-dimensional slice of integers where each index represents a word and
// its inner indexes represent the lines of each letter of the word
// For example: "Hello\nworld" will become -> [[362 623 686 686 713] [785 713 740 686 614]]
// Note that we seperate words only when we encounter "\n" and not " "
func LocateLines(input string) [][]int {
	myString := strings.ReplaceAll(input, "\\n", "\n")
	// Split into different words when encountering "\n"
	myArray := strings.Split(myString, "\n")

	// Initialize an empty slice of slices to store the lines we need to start from
	lines := make([][]int, len(myArray))

	for i, word := range myArray {
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

// WriteToTerminal writes to terminal the final output we need to print
func WriteToTerminal(linesToPrint []string) {
	for _, line := range linesToPrint {
		fmt.Println(line)
	}
}

// WriteToFile writes to a file the final output we need to print
func WriteToFile(linesToPrint []string, outputFile string) {
	file, err := os.Create(outputFile)
	if err != nil {
		//fmt.Println("Error creating file:", err)
		PrintUsage()
		return
	}
	defer file.Close()

	for _, line := range linesToPrint {
		file.WriteString(line + "\n")
	}
}

// CreateOutput combines "LocateLines" and "PrintLinesFromArray" functions to get the final output.
// This output will be printed to the terminal or written to a file with either the standard font or a specified font from the user.
func CreateOutput(input, font string) []string {
	lines := LocateLines(input)
	output, err := PrintLinesFromArray(font, lines)
	if err != nil {
		//fmt.Println("Error creating file:", err)
		PrintUsage()
	}
	return output
}

// ExtractFileName is a helper function used to extract the file name of the file we want to write the output to, if the flag "--output=" exists
func ExtractFileName(s string) string {
	return s[9:]
}
