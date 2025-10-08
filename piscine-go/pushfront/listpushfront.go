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

// ListPushFront adds a new node with the given data to the front of the list
func ListPushFront(l *List, data interface{}) {
	// If the list is empty, create a new node and set both Head and Tail to it
	if l.Head == nil {
		l.Head, l.Tail = &NodeL{Data: data}, l.Head
	} else {
		// If the list is not empty, create a new node and insert it at the front
		newNode := &NodeL{Data: data}          // Create a new node with the provided data
		newNode.Next, l.Head = l.Head, newNode // Set the Next pointer of the new node to the current Head,
		// then update the Head pointer to point to the new node
	}
}
