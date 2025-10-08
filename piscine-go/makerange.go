package piscine

func MakeRange(min, max int) []int {
	// If 'min' is greater than or equal to 'max', return nil.
	if min >= max {
		return nil
	}
	// Initialize a slice with length equal to the difference between 'max' and 'min'.
	array := make([]int, max-min)

	for i := 0; i < max-min; i++ { // Iterate through the range from 0 to the difference between 'max' and 'min'.
		array[i] = i + min // Assign the value of 'i + min' to the corresponding index in the 'array' slice.
	}
	// Return the resulting slice.
	return array
}
