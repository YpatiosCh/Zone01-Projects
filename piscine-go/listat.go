package piscine

// ListAt returns the pointer to the NodeL at the given position pos in the linked list
func ListAt(l *NodeL, pos int) *NodeL {
	iterator := l // Start from the head of the list
	inc := 0      // Initialize a counter to keep track of the current position

	// Iterate through the list until reaching the end or the desired position
	for iterator != nil {
		// If the current position matches the desired position, return the pointer to the current node
		if pos == inc {
			return iterator
		}
		inc++                    // Increment the position counter
		iterator = iterator.Next // Move to the next node
	}

	return nil // Return nil if the desired position is not found or the list is empty
}
