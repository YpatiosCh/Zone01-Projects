package piscine

// NodeI represents a node in a linked list with integer data
type NodeI struct {
	Data int    // Data stored in the node
	Next *NodeI // Pointer to the next node in the list
}

// ListSort sorts a linked list in ascending order
func ListSort(l *NodeI) *NodeI {
	Head := l // Assign the head of the list to Head

	// If the list is empty, return nil
	if Head == nil {
		return nil
	}

	// Recursively sort the list
	Head.Next = ListSort(Head.Next)

	// If the next node is not nil and the data of the current node is greater than the data of the next node,
	// rearrange the nodes
	if Head.Next != nil && Head.Data > Head.Next.Data {
		Head = move(Head)
	}

	return Head // Return the head of the sorted list
}

// move rearranges the nodes to maintain the sorting order
func move(l *NodeI) *NodeI {
	p := l      // Initialize p to the current node
	n := l.Next // Initialize n to the next node
	ret := n    // Initialize ret to n (to be returned later)

	// Move p and n until reaching a node where the data is greater than the data of the current node
	for n != nil && l.Data > n.Data {
		p = n
		n = n.Next
	}

	// Rearrange the pointers to insert the current node in the sorted position
	p.Next = l
	l.Next = n

	return ret // Return the new head of the sorted list
}
