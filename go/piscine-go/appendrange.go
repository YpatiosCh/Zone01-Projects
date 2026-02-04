package piscine

func AppendRange(min, max int) []int {
	// Initialize an empty slice to store the result.
	var array []int

	for i := min - 1; i < max-1; i++ { // Iterate through the range from 'min' to 'max'.
		array = append(array, i+1) // Append the current integer to the 'array' slice.
	}
	return array // Return the resulting slice.
}
