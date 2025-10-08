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

// ListPushBack adds a new node with the given data to the end of the list
func ListPushBack(l *List, data interface{}) {
	// Create a new node with the provided data
	n := &NodeL{Data: data}

	// If the list is empty, set both Head and Tail to the new node
	if l.Head == nil {
		l.Head = n
		l.Tail = n
	} else {
		// If the list is not empty, append the new node to the end
		l.Tail.Next = n // Set the Next pointer of the current Tail node to the new node
		l.Tail = n      // Update the Tail pointer to point to the new node
	}
}
