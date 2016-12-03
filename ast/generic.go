package ast

import (
	"fmt"
	//"psh/lex"
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

// func (r *GenericNode) Equals(node ast.Node) bool {
// 	other, ok := node.(*GenericNode)
// 	if !ok {
// 		return false
// 	}
// 	if len(r.Children) != len(other.Children) {
// 		return false
// 	}
// 	for i, _ := range r.Children {
// 		if !r.Children[i].Equals(other.Children[i]) {
// 			return false
// 		}
// 	}
// 	return t
// }
