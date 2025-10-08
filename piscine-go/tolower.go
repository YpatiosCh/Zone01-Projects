package piscine

func ToLower(s string) string {
	// Convert string to rune slice to make it mutable
	chars := []rune(s)
	// Iterate over each character
	for i, char := range chars {
		// Check if the character is uppercase
		if 'A' <= char && char <= 'Z' {
			// Convert uppercase to lowercase
			chars[i] = char + 32
		}
	}
	// Convert back to string and return
	return string(chars)
}
