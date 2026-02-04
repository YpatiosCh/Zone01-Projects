package piscine

// BTreeSearchItem searches for a node with a specific element in a binary search tree.
// It returns the node if found, or nil if the element is not in the tree.
// The function uses a recursive approach to traverse the tree, comparing the element with the current node's data.
func BTreeSearchItem(root *TreeNode, elem string) *TreeNode {
	// Check if the current node is empty or nil.
	// If it is, return nil as the element is not found in the tree.
	if root == nil {
		return nil
	}
	// If the element is less than the current node's data,
	// search the left subtree for the element.
	if elem < root.Data {
		return BTreeSearchItem(root.Left, elem)
	} else if elem > root.Data {
		// If the element is greater than the current node's data,
		// search the right subtree for the element.
		return BTreeSearchItem(root.Right, elem)
	} else {
		// If the element matches the current node's data,
		// return the current node as the element is found.
		return root
	}
}
