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

// ListLast returns the data interface of the last element of the linked list
func ListLast(l *List) interface{} {
	// If the list is empty, return nil
	if l.Tail == nil {
		return nil
	}

	return l.Tail.Data // Return the data stored in the Tail node
}
