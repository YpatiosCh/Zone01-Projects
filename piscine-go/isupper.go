package piscine

func IsUpper(s string) bool {
	// Iterate over each character in the string
	for _, char := range s {
		// Check if the character is not an uppercase letter
		if !('A' <= char && char <= 'Z') {
			return false
		}
	}
	// If all characters are uppercase, return true
	return true
}
