package piscine

func IsPrintable(s string) bool {
	// Iterate over each character in the string
	for _, char := range s {
		// Check if the character is not printable
		if char < 32 || char > 126 {
			return false
		}
	}
	// If all characters are printable, return true
	return true
}
