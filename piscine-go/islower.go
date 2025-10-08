package piscine

func IsLower(s string) bool {
	// Iterate over each character in the string
	for _, char := range s {
		// Check if the character is not an lowerrcase letter
		if !('a' <= char && char <= 'z') {
			return false
		}
	}
	// If all characters are lowercase, return true
	return true
}
