package piscine

// NodeL represents a node in a linked list
type NodeL struct {
	Data interface{} // Data stored in the node
	Next *NodeL      // Pointer to the next node in the list
}

// List represents a linked list
type List struct {
	Head *NodeL // Pointer to the first node in the list
	Tail *NodeL // Pointer to the last node in the list
}

// ListSize returns the number of elements in the linked list
func ListSize(l *List) int {
	size := 0         // Initialize a variable to store the size of the list
	current := l.Head // Start from the head of the list

	// Traverse the list until reaching the end (current == nil)
	for current != nil {
		size++                 // Increment the size counter
		current = current.Next // Move to the next node
	}

	return size // Return the size of the list
}
