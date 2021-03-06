package ast

import (
	"fmt"
)

type Parselet interface {
	Parse(*Parser) error
}

type Node interface {
	fmt.Formatter
}

type Command interface {
	Node
	Parselet
	IsCommand()
}

type Expr interface {
	Node
	Parselet
	IsExpr()
}

type StrPiece interface {
	fmt.Formatter
	IsStrPiece()
}
