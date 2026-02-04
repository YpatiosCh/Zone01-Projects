package piscine

// BTreeApplyByLevel applies a given function to each node in the binary tree level by level.
// It uses a queue to keep track of nodes to be processed, starting with the root node.
// The function is useful for performing operations on each node of the tree in a breadth-first manner.
func BTreeApplyByLevel(root *TreeNode, f func(...interface{}) (int, error)) {
	// Check if the root node is empty or nil.
	// If it is, return immediately as there's nothing to do.
	if root == nil {
		return
	}
	// Initialize a queue with the root node.
	// The queue will be used to process nodes level by level.
	queue := []*TreeNode{root}
	// Process nodes in the queue until it is empty.
	for len(queue) > 0 {
		// Dequeue the first node from the queue.
		node := queue[0]
		// Remove the dequeued node from the queue.
		queue = queue[1:]
		// Apply the given function to the data part of the dequeued node.
		// This is where the actual operation on the node's data is performed.
		f(node.Data)
		// If the dequeued node has a left child, enqueue it to the queue.
		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		// If the dequeued node has a right child, enqueue it to the queue.
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
}
