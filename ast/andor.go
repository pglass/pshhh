package ast

import (
	"fmt"

	"github.com/pglass/pshhh/lex"
)

type AndOrClause struct {
	Left        Node
	Operator    *lex.Token
	AndOrClause *AndOrClause
}

func (a *AndOrClause) IsExpr() {}

func NewAndOrClause(left Node) *AndOrClause {
	return &AndOrClause{
		Left:        left,
		Operator:    nil,
		AndOrClause: nil,
	}
}

func (a *AndOrClause) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "AndOrClause[")
	if a.Left != nil {
		fmt.Fprintf(f, "%v", a.Left)
	}
	if a.Operator != nil {
		fmt.Fprintf(f, " %v", a.Operator.Text)
	}
	if a.AndOrClause != nil {
		fmt.Fprintf(f, " %v", a.AndOrClause)
	}
	fmt.Fprintf(f, "]")
}

func (a *AndOrClause) Parse(parser *Parser) error {
	if !parser.Lexer.HasAnyToken(lex.AndIf, lex.OrIf) {
		return nil
	}

	tok := parser.Lexer.Next()
	a.Operator = &tok

	parser.ConsumeWhile(lex.Space)

	if node, err := parser.ParseNext(); err != nil {
		return err
	} else if node != nil {
		// todo: maybe only certain kinds of nodes are valid here
		a.AndOrClause = NewAndOrClause(node)
		if err := a.AndOrClause.Parse(parser); err != nil {
			return err
		}
	}

	return nil
}
