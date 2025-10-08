package piscine

func ActiveBits(n int) int {
	// Initialize a counter to keep track of the number of active bits
	count := 0

	// Iterate through each bit of the number
	for n != 0 {
		// Use bitwise AND to check if the current bit is set (1)
		// If it is, increment the counter
		count += n & 1

		// Right shift the number to move to the next bit
		n >>= 1
	}

	// Return the total number of active bits
	return count
}
