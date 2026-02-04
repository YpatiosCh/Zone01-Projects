package piscine

func SplitWhiteSpaces(str string) []string {
	index := 0
	count := 0
	word := ""

	// Count the number of substrings by checking consecutive white space characters.
	for i, v := range str {
		if v == ' ' && str[i+1] != ' ' {
			count++
		}
	}
	// Initialize a slice to store the resulting substrings.
	result := make([]string, count+1)
	// Iterate over the input string 'str'.
	for _, r := range str {
		// Check if the current character is a separator.
		if isSeparator(r) {
			// If 'word' is not empty, add it to the 'result' slice and reset 'word'.
			if word != "" {
				result[index] = word
				index++
				word = ""
			}
		} else {
			// If the current character is not a separator, add it to the 'word'.
			word += string(r)
		}
	}
	// Add the remaining 'word' to the 'result' slice if it's not empty.
	size := 0
	for z := range result {
		size++
		z++
	}
	if word != "" {
		result[size-1] = word
	}
	// Return the resulting slice.
	return result
}

// isSeparator checks if the given rune 'r' is a white space character.
func isSeparator(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}
