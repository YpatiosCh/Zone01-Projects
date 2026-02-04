package tools

import (
	"os"
	"fmt"
	"bufio"
)

// Get ASCII art for a specific letter from the font file
func GetAsciiArtForLetter(file *os.File, letter rune) ([]string, error) {
	// Only process printable characters (ASCII values between 32 and 126)
	if letter < 32 || letter > 126 {
		return nil, fmt.Errorf("character %c is not printable", letter)
	}

	line := int(letter)
	lineFromText := (line-32)*9 + 1

	_, err := file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	// Skip to the correct starting line for the letter in the font file
	for i := 0; i < lineFromText; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to find line %d for letter %c", lineFromText, letter)
		}
	}

	// Read the 8 lines of ASCII art for the character
	var asciiArt []string
	for i := 0; i < 8; i++ {
		if scanner.Scan() {
			asciiArt = append(asciiArt, scanner.Text())
		} else {
			return nil, fmt.Errorf("failed to read full ascii art for letter %c (expected 8 lines, got %d)", letter, i)
		}
	}

	return asciiArt, nil
}