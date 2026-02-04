package piscine

// ListReverse reverses the order of elements in the linked list
func ListReverse(l *List) {
	current := l.Head     // Initialize current node to the head of the list
	var prev *NodeL = nil // Initialize prev node to nil

	// Traverse the list and reverse the pointers
	for current != nil {
		next := current.Next // Store the next node
		current.Next = prev  // Reverse the pointer of current node to point to the previous node
		prev = current       // Move prev to current node
		current = next       // Move current to the next node
	}

	// Swap the Head and Tail pointers to complete the reversal
	temp := l.Head
	l.Head = l.Tail
	l.Tail = temp
}
