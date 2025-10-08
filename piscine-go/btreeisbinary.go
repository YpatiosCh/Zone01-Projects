package piscine

// BTreeIsBinary checks if a binary tree is a binary search tree (BST).
// A BST is a tree where for each node, all elements in the left subtree are less than the node,
// and all elements in the right subtree are greater than the node.
// This function uses a recursive approach to traverse the tree and validate the BST property.
func BTreeIsBinary(root *TreeNode) bool {
	// If the current node is empty or nil, it is considered a BST.
	if root == nil {
		return true
	}
	// If the current node has both left and right children,
	// check if the left child's data is less than the current node's data
	// and if the right child's data is greater than the current node's data.
	// If either condition is not met, return false.
	if root.Left != nil && root.Right != nil {
		if root.Left.Data > root.Data || root.Right.Data < root.Data {
			return false
		}
		// Recursively check if the left and right subtrees are also BSTs.
		// If either subtree is not a BST, return false.
		return BTreeIsBinary(root.Left) && BTreeIsBinary(root.Right)
	}
	// If the current node has only one child or no children, it is considered a BST.
	return true
}
