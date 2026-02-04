package piscine

func ToUpper(s string) string {
	// Convert string to rune slice to make it mutable
	chars := []rune(s)
	// Iterate over each character
	for i, char := range chars {
		// Check if the character is lowercase
		if 'a' <= char && char <= 'z' {
			// Convert lowercase to uppercase
			chars[i] = char - 32
		}
	}
	// Convert back to string and return
	return string(chars)
}
