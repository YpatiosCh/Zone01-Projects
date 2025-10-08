package piscine

import "github.com/01-edu/z01"

// PrintWordsTables prints each string in the 'table' slice, separated by newlines.
func PrintWordsTables(table []string) {
	// Iterate over each string in the 'table' slice.
	for _, v := range table {
		// Call the PRune function to print each string rune by rune.
		PRune(v)
	}
}

// PRune prints each rune of the input string 'str' followed by a newline.
func PRune(str string) {
	strRune := []rune(str)
	// Iterate over each rune in the string 'str'.
	for _, v := range strRune {
		// Print the current rune.
		z01.PrintRune(v)
	}
	// Print a newline character after printing all the runes of the string.
	z01.PrintRune('\n')
}
