package ascii

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Function to print ASCII art for a word or a whole index of a []sting
func PrintAsciiArt(word, filename, color, outputFile, alignFlag string, colorFlags []bool) error {

	// Open the ASCII art file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("\nfont not available. You can choose between: <standard> <shadow> <thinkertoy>")
	}
	defer file.Close()
	// Prepare lines for ASCII art
	lines := make([]string, 8) // Each letter has 8 lines in the ASCII art

	var length int
	// to count how long is the width of the ascii fullInput
	for _, letter := range word {
		_, index, _ := GetAsciiArtForLetter(file, letter)
		length += index
	}
	// Iterate over each letter in the word
	for i, letter := range word {
		art, _, err := GetAsciiArtForLetter(file, letter)
		if err != nil {
			return fmt.Errorf("\nerror reading ascii art for letter %c: %v", letter, err)
		}

		// Use color based on the boolean flag
		currentColor := Reset // Default to reset (white)
		if colorFlags[i] {
			currentColor = color
		}

		if alignFlag == "justify" && isEmpty(art) {
			Justify(art, length, word)
		}

		// Append each line of the letter directly to the corresponding line in 'lines'
		for j := 0; j < 8; j++ {
			lines[j] += fmt.Sprintf("%s%s%s", currentColor, art[j], Reset) // Append with the specified color
		}
	}
	termWidth, err := GetTerminalWidth()
	if err != nil {
		PrintUsage()
		fmt.Println(err)
	}

	if outputFile != "" {
		WriteToFile(outputFile, lines)
	} else {
		// Apply selected alignment
		switch alignFlag {
		case "center":
			AlignCenter(lines, termWidth, length)
		case "right":
			AlignRight(lines, termWidth, length)
		default:
			// default will be on the left or jusify if it was specified in the flag
			// Print to terminal
			for _, line := range lines {
				fmt.Println(line)
			}
		}
	}

	return nil
}

// helper function for PrintAsciiArt
// GetAsciiArtForLetter reads ASCII art for a specific character from the file
func GetAsciiArtForLetter(file *os.File, letter rune) ([]string, int, error) {
	line := int(letter)
	lineFromText := (line-32)*9 + 1 // Adjusted formula to find the starting line
	// Move to the appropriate line in the file
	_, err := file.Seek(0, 0) // Start from the beginning of the file
	if err != nil {
		return nil, 0, err
	}

	scanner := bufio.NewScanner(file)

	// Skip lines until we reach the starting line for the letter
	for i := 0; i < lineFromText; i++ {
		if !scanner.Scan() {
			return nil, 0, fmt.Errorf("\nfailed to find line %d for letter %c", lineFromText, letter)
		}
	}

	// Read the next 8 lines to get the ASCII art for the letter
	var asciiArt []string
	for i := 0; i < 8; i++ {
		if scanner.Scan() {
			asciiArt = append(asciiArt, scanner.Text())
		} else {
			return nil, 0, fmt.Errorf("\nfailed to read full ascii art for letter %c (expected 8 lines, got %d)", letter, i)
		}
	}
	index := len(asciiArt[0])
	return asciiArt, index, nil
}

// PrintFullAscii prints the whole ascii-art with the user input and correctly handling the "\n" if there is one
func PrintFullAscii(substring, fullInput, colorFlag, outputFile, font, alignFlag string) {
	// Split the input on '\n' (newline) to handle multi-line input
	parts := strings.Split(fullInput, "\\n")

	color := determineColor(&colorFlag)

	// Process each part and determine which characters should be colored
	for _, part := range parts {
		if part == "" {
			fmt.Println()
			continue
		}

		// Create color flags for the current part
		colorSubs := determineSubToColor(part, substring)

		// Print the ASCII art for the part with the appropriate colors
		err := PrintAsciiArt(part, font, color, outputFile, alignFlag, colorSubs)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// IsSubstring function determines whether a string is a substring of another string
func IsSubstring(str, subStr string) bool {
	return strings.Contains(str, subStr)
}

// PrintUsage prints to the terminal the correct way to use the program
func PrintUsage() {
	fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
	fmt.Println("\nEX: go run . --align=right something standard")
}

// IsFont returns true if the string we check is a fontName
func IsFont(s string) bool {
	return s == "standard" || s == "shadow" || s == "thinkertoy"
}

func isEmpty(art []string) bool {
	for _, line := range art {
		for _, r := range line {
			if r != ' ' {
				return false
			}
		}
	}
	return true
}
