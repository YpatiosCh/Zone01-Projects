package piscine

// IsPositiveNode checks if the data of the given node is positive
func IsPositiveNode(node *NodeL) bool {
	switch node.Data.(type) {
	case int, float32, float64, byte:
		return node.Data.(int) > 0
	default:
		return false
	}
}

// IsAlNode checks if the data of the given node is not numeric
func IsAlNode(node *NodeL) bool {
	switch node.Data.(type) {
	case int, float32, float64, byte:
		return false
	default:
		return true
	}
}

// ListForEachIf applies the function f to each node in the linked list l that satisfies the condition specified by the function cond
func ListForEachIf(l *List, f func(*NodeL), cond func(*NodeL) bool) {
	// Start from the head of the list
	current := l.Head

	// Iterate through each node in the list
	for current != nil {
		// If the condition specified by the function cond is true for the current node, apply the function f
		if cond(current) {
			f(current)
		}
		// Move to the next node
		current = current.Next
	}
}
