package piscine

// ListRemoveIf removes nodes from the linked list l that contain the specified data_ref
func ListRemoveIf(l *List, data_ref interface{}) {
	temp := l.Head // Initialize temp to the head of the list
	prev := l.Head // Initialize prev to the head of the list

	// Remove nodes from the beginning of the list that contain data_ref
	for temp != nil && temp.Data == data_ref {
		l.Head = temp.Next // Update the head of the list to the next node
		temp = l.Head      // Move temp to the new head of the list
	}

	// Iterate through the list to remove nodes containing data_ref
	for temp != nil {
		// Find the node containing data_ref
		for temp != nil && temp.Data != data_ref {
			prev = temp      // Update prev to the current node
			temp = temp.Next // Move temp to the next node
		}

		// If temp is nil, data_ref is not found in the remaining nodes
		if temp == nil {
			return
		}

		// Remove the node containing data_ref
		prev.Next = temp.Next // Update the next pointer of the previous node to skip temp
		temp = prev.Next      // Move temp to the next node
	}
}
