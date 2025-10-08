package piscine

// SortedListMerge merges two sorted linked lists l1 and l2 into a single sorted linked list
func SortedListMerge(l1 *NodeI, l2 *NodeI) *NodeI {
	// Sort the input lists l1 and l2
	l1 = ListSort(l1)
	l2 = ListSort(l2)

	// If l1 is empty, return l2
	if l1 == nil {
		return l2
	}

	// If l2 is empty, return l1
	if l2 == nil {
		return l1
	}

	// Merge the sorted lists
	if l1.Data <= l2.Data {
		// If the data of the head node of l1 is less than or equal to the data of the head node of l2,
		// recursively merge l1.Next and l2, and set the result as the next node of l1
		l1.Next = SortedListMerge(l1.Next, l2)
		return l1 // Return l1 as the merged list
	} else {
		// If the data of the head node of l2 is less than the data of the head node of l1,
		// recursively merge l1 and l2.Next, and set the result as the next node of l2
		l2.Next = SortedListMerge(l1, l2.Next)
		return l2 // Return l2 as the merged list
	}
}
