package piscine

// CompStr compares two interfaces and returns true if they are equal, false otherwise
func CompStr(a, b interface{}) bool {
	return a == b
}

// ListFind searches for a reference element in the linked list l using the comparison function comp
func ListFind(l *List, ref interface{}, comp func(a, b interface{}) bool) *interface{} {
	iterator := l.Head // Start from the head of the list

	// Iterate through each node in the list
	for iterator != nil {
		// Compare the data of the current node with the reference using the comparison function comp
		if comp(iterator.Data, ref) {
			return &iterator.Data // Return a pointer to the data if it matches the reference
		}
		iterator = iterator.Next // Move to the next node
	}
	return nil // Return nil if the reference element is not found in the list
}
