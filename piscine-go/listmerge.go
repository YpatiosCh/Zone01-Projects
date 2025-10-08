package piscine

// ListMerge merges two linked lists l1 and l2
func ListMerge(l1 *List, l2 *List) {
	// If either of the lists is nil, return
	if l2 == nil || l1 == nil {
		return
	}

	// If the head of l1 is nil, set it to the head of l2
	if l1.Head == nil {
		l1.Head = l2.Head
		l1.Tail = l2.Head // Update the tail of l1 to the tail of l2
		return
	}

	// Append the nodes of l2 to the end of l1
	l1.Tail.Next = l2.Head // Set the next pointer of the tail of l1 to the head of l2
}
