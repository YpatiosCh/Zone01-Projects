package piscine

func AtoiBase(s string, base string) int {
	// Check if the base is valid
	if len(base) < 2 {
		return 0
	}
	for i := 0; i < len(base); i++ {
		if base[i] == '+' || base[i] == '-' {
			return 0
		}
		for j := i + 1; j < len(base); j++ {
			if base[i] == base[j] {
				return 0
			}
		}
	}

	// Create a map to store the index of each digit in the base
	baseMap := make(map[rune]int)
	for i, digit := range base {
		baseMap[digit] = i
	}

	// Initialize the result
	result := 0

	// Iterate over each digit in the string
	for _, digit := range s {
		// Check if the digit is valid
		if _, ok := baseMap[digit]; !ok {
			return 0
		}
		// Update the result
		result = result*len(base) + baseMap[digit]
	}

	return result
}
