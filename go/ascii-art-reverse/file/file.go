package file

import (
	"bufio"
	"fmt"
	"os"
	"reverse/args"
	"sort"
)

// ScanFile will scan the fileName provided in the arguments.
// every line will be appended to a []string and when done it
// will return this []string
// Basically, it returns the ascii- art from the file provided.
// If there is an error while opening the file it will return the error.
func ScanFile(filename *string) ([]string, error) {
	file, err := os.Open(*filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	var asciiArt []string

	for scanner.Scan() {
		line := scanner.Text()
		asciiArt = append(asciiArt, line)
	}

	return asciiArt, nil
}

// IndexOfEachLetter gets the ascii-art from the text file provided
// and returns a []int where each number represents the starting and the ending point of the letter
// in ascii form.
// We need the starting and ending point to split the word in letters and for each letter make a deep
// comparison in the standard.txt file to determine which ascii character it is.
// This function uses getIndexes and handleSpaces helper functions.
func IndexOfEachLetter(asciiArt []string) []int {
	indexes := GetIndexes(asciiArt)

	handledSpaces := handleSpaces(indexes)

	return handledSpaces
}

// getIndexes finds the start/end indexes of each ascii-art letter.
// It will helps us store each ascii-art letter seperately,
// so that we can make a deep comparison later in the standard.txt file.
func GetIndexes(asciiArt []string) []int {
	var result []int
	iMap := make(map[int]int)

	for _, line := range asciiArt {
		for i, char := range line {
			if char == ' ' {
				iMap[i]++
			}
		}
	}

	for key, value := range iMap {
		if value == 8 {
			result = append(result, key)
		}
	}

	sort.Ints(result)

	return result
}

// helper function of IndexOfEachLetter.
// Since the ascii-art representation of space (" ") is all spaces,
// we need to determine where the space starts and ends.
// handleSpaces function handles this situation and returns
// the final []int that tells us where each letter starts and ends.
func handleSpaces(indexes []int) []int {
	var result []int
	result = append(result, 0)

	for i, num := range indexes {
		if i != len(indexes)-1 && indexes[i+1] == num+1 && indexes[i-1] == num-1 {
			continue
		} else {
			result = append(result, num+1)
		}
	}
	return result

}

// PrintLetter helps us seperate each ascii-art letter separately,
// using start and end indexes of each letter.
func PrintLetter(start, end int, asciiArt []string) []string {
	var result []string

	for _, line := range asciiArt {
		result = append(result, line[start:end])
	}
	return result
}

// SearchFile function takes the separate ascii-art letter and
// makes a deep comparison in the standard.txt.
// When we find 8 consecutive matches, we return the first number of line
// that we encountered the match.
// This line will be used to determine wich character of the ascii table we want to print.
// If we encounter an error, we will return it.
func SearchFile(art []string) (int, error) {
	file, err := os.Open("./fonts/standard.txt")
	if err != nil {
		args.PrintReverseUsage()
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	consecutiveMatchCount := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Check if current line matches the next expected line in art
		if line == art[consecutiveMatchCount] {
			consecutiveMatchCount++

			// Check if we've found all 8 consecutive matches
			if consecutiveMatchCount == len(art) {
				return lineNumber - len(art) + 1, nil
			}
		} else {
			// Reset if line doesn't match
			consecutiveMatchCount = 0

			// Special case: restart matching if first line matches
			if line == art[0] {
				consecutiveMatchCount = 1
			}
		}
	}

	return -1, nil
}

// GetLetter uses the below mathematical formula to convert the line number we found from
// SearchFile function to the corresponding rune.
func GetLetter(line int) rune {
	letter := ((line - 1) / 9) + 32
	return rune(letter)
}

// PrintArtToAscii gets the starting/ending indexes of each ascii-art letter
// and the ascii-art itself in order to seperate each ascii-art letter.
// It then makes a deep comparison in the standard.txt file using SearchFile function.
// When we find the match, we convert the number of line in the text where the match was made,
// convert it to rune and add it to string result as the final output.
func PrintArtToAscii(indexes []int, asciiArt []string) {
	var result string
	for i := 0; i < len(indexes)-1; i++ {
		start := indexes[i]
		end := indexes[i+1]
		art := PrintLetter(start, end, asciiArt)
		line, err := SearchFile(art)
		if err != nil {
			args.PrintReverseUsage()
			return
		}
		letter := GetLetter(line)
		result += string(letter)
	}
	fmt.Println(result)
}
