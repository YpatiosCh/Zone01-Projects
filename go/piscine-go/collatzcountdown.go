package piscine

// CollatzCountdown calculates the number of steps it takes for a number to reach 1 using the Collatz conjecture.
// If the starting number is 0 or negative, it returns -1.
func CollatzCountdown(start int) int {
	count := 0 // Initialize steps counter to 1, assuming the starting number itself counts as a step
	if start <= 0 {
		return -1 // Return -1 for invalid input
	} else {
		// Continue iterating until start becomes 1
		for start != 1 {
			// Apply Collatz conjecture
			if start%2 == 0 {
				start = start / 2
			} else if start%2 == 1 {
				start = start*3 + 1
			}
			count++ // Increment the steps counter
		}
	}
	return count // Return the total steps taken
}
