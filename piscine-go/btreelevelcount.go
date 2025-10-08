package piscine

// BTreeLevelCount calculates the number of levels in a binary tree.
// It uses a recursive approach to traverse the tree and count the levels.
// The function returns the number of levels in the tree, which is determined by the depth of the tree.
func BTreeLevelCount(root *TreeNode) int {
	// Check if the current node is empty or nil.
	// If it is, return 0 as the tree has no levels.
	if root == nil {
		return 0
	}
	// Recursively calculate the number of levels in the left subtree.
	left := BTreeLevelCount(root.Left)
	// Recursively calculate the number of levels in the right subtree.
	right := BTreeLevelCount(root.Right)
	// Determine the maximum number of levels between the left and right subtrees.
	// Add 1 to account for the current level (the root node).
	if left > right {
		return left + 1
	}
	return right + 1
}
