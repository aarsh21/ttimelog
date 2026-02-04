package treeview

import "strings"

type TreeNode struct {
	Label    string
	Children []*TreeNode
	Expanded bool
	Path     string
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

func AppendPath(rootNode *TreeNode, path []string, index int) {
	// Base case: no more labels to consume
	if len(path) == index {
		return
	}

	currentLabel := path[index]

	for _, child := range rootNode.Children {
		if child.Label == currentLabel {
			AppendPath(child, path, index+1)
			return
		}
	}

	// TODO: improve this to only save "Path" for leaf nodes
	newChild := &TreeNode{Label: currentLabel, Path: strings.Join(path, ":")}
	rootNode.Children = append(rootNode.Children, newChild)

	AppendPath(newChild, path, index+1)
}
