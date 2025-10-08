package piscine

// TreeNode represents a node in a binary tree.
// It contains pointers to its left and right children, its parent, and the data it holds.
type TreeNode struct {
	Left, Right, Parent *TreeNode
	Data                string
}

// BTreeInsertData inserts a new node with the given data into the binary tree.
// If the tree is empty, it creates a new root node.
// Otherwise, it recursively finds the correct position for the new node based on the data comparison.
func BTreeInsertData(root *TreeNode, data string) *TreeNode {
	// If the root is nil, it means the tree is empty.
	// Create a new node with the given data and return it as the new root.
	if root == nil {
		return &TreeNode{nil, nil, nil, data}
	}
	// If the data to be inserted is less than the root's data,
	// insert the data into the left subtree.
	if data < root.Data {
		root.Left = BTreeInsertData(root.Left, data)
		// Set the parent of the newly inserted node to the current root.
		root.Left.Parent = root
	} else {
		// If the data to be inserted is greater than or equal to the root's data,
		// insert the data into the right subtree.
		root.Right = BTreeInsertData(root.Right, data)
		// Set the parent of the newly inserted node to the current root.
		root.Right.Parent = root
	}
	// Return the root node, which may have been updated if a new node was inserted.
	return root
}
