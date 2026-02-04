package tools

import (
	"os"
	"fmt"
	"strings"
)

// Print ASCII art for the given text and banner
func PrintAsciiArt(word string, filename string) ([]string, error) {
	// Open the font file for the selected banner
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("font not available. You can choose between: standard, shadow, thinkertoy")
	}
	defer file.Close()

	// Split the input text into words (if there are multiple words)
	words := strings.Split(strings.ReplaceAll(word, "\r\n", "\n"), "\n")
	var allLines []string

	// Loop through each word and generate ASCII art for it
	for wordIndex, currentWord := range words {
		if currentWord == "" {
			if wordIndex < len(words)-1 {
				allLines = append(allLines, "")
			}
			continue
		}

		lines := make([]string, 8) // Each letter has 8 lines in the ASCII art

		// Loop through each letter in the current word and generate ASCII art
		for _, letter := range currentWord {
			art, err := GetAsciiArtForLetter(file, letter)
			if err != nil {
				return nil, fmt.Errorf("error reading ascii art for letter %c: %v", letter, err)
			}
			// Append the ASCII art for each letter to the lines
			for j := 0; j < 8; j++ {
				lines[j] += art[j]
			}
		}

		allLines = append(allLines, lines...)
		if wordIndex < len(words)-1 {
			allLines = append(allLines, "")
		}
	}

	return allLines, nil
}