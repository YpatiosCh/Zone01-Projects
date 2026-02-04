package piscine

func Abort(a, b, c, d, e int) int {
	// Create a slice to hold the arguments
	numbers := []int{a, b, c, d, e}

	// Sort the slice
	for i := 0; i < len(numbers)-1; i++ {
		for j := i + 1; j < len(numbers); j++ {
			if numbers[i] > numbers[j] {
				numbers[i], numbers[j] = numbers[j], numbers[i]
			}
		}
	}

	// Return the median (middle value)
	return numbers[2] // Since the slice is sorted, the median is the third element
}
