package ast

import (
	"fmt"
)

// The root node represents any shell code
type GenericNode struct {
	Children []Node
}

func NewGenericNode(children ...Node) *GenericNode {
	return &GenericNode{Children: children}
}

func (r *GenericNode) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "GenericNode[")
	if len(r.Children) > 0 {
		fmt.Fprintf(f, "%v", r.Children[0])
		for _, child := range r.Children[1:] {
			fmt.Fprintf(f, "\n%v", child)
		}
	}
	fmt.Fprintf(f, "]")
}
