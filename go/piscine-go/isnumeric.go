package piscine

func IsNumeric(s string) bool {
	// Iterate over each character in the string
	for _, char := range s {
		// Check if the character is not numeric
		if !('0' <= char && char <= '9') {
			return false
		}
	}
	// If all characters are numeric or the string is empty, return true
	return true
}
