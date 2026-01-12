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
