package piscine

func IsAlpha(s string) bool {
	// Iterate over each character in the string
	for _, char := range s {
		// Check if the character is not alphanumeric
		if !('0' <= char && char <= '9') && !('a' <= char && char <= 'z') && !('A' <= char && char <= 'Z') {
			return false
		}
	}
	// If all characters are alphanumeric or the string is empty, return true
	return true
}
