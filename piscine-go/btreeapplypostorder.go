package piscine

// BTreeApplyPostorder applies a given function to each node in the binary tree in a post-order traversal.
// Post-order traversal visits the left subtree, then the right subtree, and finally the root node.
// This function is useful for performing operations on each node of the tree in a specific order,
// often used for tasks that require processing the tree's nodes in a bottom-up manner.
func BTreeApplyPostorder(root *TreeNode, f func(...interface{}) (int, error)) {
	// Check if the current node is empty or nil.
	// If it is, return immediately as there's nothing to do.
	if root == nil {
		return
	}
	// Traverse the left subtree by recursively calling the post-order function.
	// This ensures that all nodes in the left subtree are visited before the root.
	BTreeApplyPostorder(root.Left, f)
	// Traverse the right subtree by recursively calling the post-order function.
	// This ensures that all nodes in the right subtree are visited before the root.
	BTreeApplyPostorder(root.Right, f)
	// Apply the given function to the data part of the root (or current node).
	// This is where the actual operation on the node's data is performed.
	f(root.Data)
}
