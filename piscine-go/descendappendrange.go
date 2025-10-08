package piscine

func DescendAppendRange(max, min int) []int {
	// Initialize an empty slice of integers
	var result []int

	// Check if max is less than or equal to min
	if max <= min {
		return []int{} // Return an empty slice explicitly
	}

	// Iterate from max down to min (exclusive)
	for i := max; i > min; i-- {
		// Append each value to the result slice
		result = append(result, i)
	}

	return result
}
