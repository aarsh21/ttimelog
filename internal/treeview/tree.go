package treeview

type TreeNode struct {
	Label    string
	Children []*TreeNode
	Expanded bool
}

type Row struct {
	TreeNode *TreeNode
	Depth    int
}

func Traverse(node *TreeNode, depth int, rows *[]Row) {
	if node == nil {
		return
	}

	*rows = append(*rows, Row{
		TreeNode: node,
		Depth:    depth,
	})

	if !node.Expanded {
		return
	}

	for _, child := range node.Children {
		Traverse(child, depth+1, rows)
	}
}

func AppendPath(rootNode *TreeNode, path []string) {
	// Base case: no more labels to consume
	if len(path) == 0 {
		return
	}

	currentLabel := path[0]

	for _, child := range rootNode.Children {
		if child.Label == currentLabel {
			AppendPath(child, path[1:])
			return
		}
	}

	newChild := &TreeNode{Label: currentLabel}
	rootNode.Children = append(rootNode.Children, newChild)

	AppendPath(newChild, path[1:])
}
