package ast

import (
	"fmt"

	"github.com/pglass/pshhh/lex"
)

type DoClause struct {
	Children []Node
}

func NewDoClause() *DoClause {
	return &DoClause{Children: []Node{}}
}

func (d *DoClause) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "do %v done", d.Children)
}

func (d *DoClause) Parse(parser *Parser) error {
	// "do"
	if _, err := parser.ConsumeToken(lex.Do, nil); err != nil {
		return err
	}

	for {
		parser.ConsumeWhile(lex.Space, lex.Newline)

		// stop when we see "done"
		if parser.Lexer.HasAnyToken(lex.Done) {
			parser.Lexer.Next()
			break
		}

		parser.ConsumeWhile(lex.Space, lex.Newline)

		if node, err := parser.ParseNext(); node != nil {
			// todo: only certain node types are valid here?
			d.Children = append(d.Children, node)
		} else if err != nil {
			return err
		}
	}

	if len(d.Children) == 0 {
		return fmt.Errorf("Empty do...done block")
	}

	return nil
}
