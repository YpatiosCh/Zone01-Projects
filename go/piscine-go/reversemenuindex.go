package piscine

func ReverseMenuIndex(menu []string) []string {
	// Determine the length of the input slice
	menuLen := len(menu)
	// Create a new slice with the same length as the input slice
	output := make([]string, menuLen)
	// Iterate over the input slice
	for i, n := range menu {
		// Calculate the index in the output slice where the current element should be placed
		// This is done by subtracting the current index from the total length of the menu
		j := menuLen - i - 1

		// Assign the current element from the input slice to the calculated index in the output slice
		output[j] = n
	}
	// Return the reversed slice
	return output
}
