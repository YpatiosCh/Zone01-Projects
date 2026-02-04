package piscine

// SortListInsert inserts a node with data_ref into a sorted linked list l
func SortListInsert(l *NodeI, data_ref int) *NodeI {
	// Create a new node with the data_ref
	n := &NodeI{}
	n.Data = data_ref
	n.Next = nil

	// If the list is empty or the data of the head node is greater than or equal to the data_ref,
	// insert the new node at the beginning of the list
	if l == nil || l.Data >= n.Data {
		n.Next = l // Connect the new node to the current head
		return n   // Return the new node as the new head of the list
	} else {
		temp := l // Initialize a temporary pointer to traverse the list
		// Traverse the list until finding the correct position to insert the new node
		for temp.Next != nil && temp.Next.Data < n.Data {
			temp = temp.Next
		}
		// Insert the new node between temp and temp.Next
		n.Next = temp.Next // Connect the new node to the next node of temp
		temp.Next = n      // Connect temp to the new node
	}

	return l // Return the head of the list
}
