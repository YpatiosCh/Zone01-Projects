package piscine

func StringToIntSlice(str string) []int {
	// Initialize an empty slice of integers
	var ints []int

	// Iterate over each rune (character) in the string
	for _, r := range str {
		// Convert the rune to its ASCII value
		intValue := int(r)

		// Append the integer value to the slice
		ints = append(ints, intValue)
	}

	return ints
}
