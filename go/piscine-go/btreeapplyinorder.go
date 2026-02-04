package piscine

// BTreeApplyInorder applies a given function to each node in the binary tree in an in-order traversal.
// In-order traversal visits the left subtree, then the root node, and finally the right subtree.
// This function is useful for performing operations on each node of the tree in a sorted order.
func BTreeApplyInorder(root *TreeNode, f func(...interface{}) (int, error)) {
	// Check if the current node is empty or nil.
	// If it is, return immediately as there's nothing to do.
	if root == nil {
		return
	}
	// Traverse the left subtree by recursively calling the in-order function.
	// This ensures that all nodes in the left subtree are visited before the root.
	BTreeApplyInorder(root.Left, f)
	// Apply the given function to the data part of the root (or current node).
	// This is where the actual operation on the node's data is performed.
	f(root.Data)
	// Traverse the right subtree by recursively calling the in-order function.
	// This ensures that all nodes in the right subtree are visited after the root.
	BTreeApplyInorder(root.Right, f)
}
