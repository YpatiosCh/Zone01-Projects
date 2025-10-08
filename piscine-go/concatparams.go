package piscine

func ConcatParams(args []string) string {
	var result string
	size := 0

	// Calculate the size of the 'args' slice.
	for i := range args {
		i++    // Incrementing 'i' to iterate over each element.
		size++ // Incrementing 'size' to count the number of elements.
	}
	// Iterate through the 'args' slice.
	for i, v := range args {
		// Concatenate the current element to the 'result' string.
		result += v
		// Add a newline separator if it's not the last element.
		if i != size-1 {
			result += "\n"
		}
	}
	// Return the concatenated string.
	return result
}
