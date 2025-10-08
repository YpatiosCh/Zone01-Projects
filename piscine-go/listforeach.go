package piscine

// ListForEach applies the function f to each node in the linked list l
func ListForEach(l *List, f func(*NodeL)) {
	// Start from the head of the list
	current := l.Head

	// Iterate through each node in the list
	for current != nil {
		// Apply the function f to the current node
		f(current)

		// Move to the next node
		current = current.Next
	}
}

// Add2_node adds 2 to the data of the given node
func Add2_node(node *NodeL) {
	switch node.Data.(type) {
	case int:
		node.Data = node.Data.(int) + 2
	case string:
		node.Data = node.Data.(string) + "2"
	}
}

// Subtract3_node subtracts 3 from the data of the given node
func Subtract3_node(node *NodeL) {
	switch node.Data.(type) {
	case int:
		node.Data = node.Data.(int) - 3
	case string:
		node.Data = node.Data.(string) + "-3"
	}
}
