package ast

import (
	"fmt"
	"psh/lex"
)

type StringPiece interface {
	fmt.Formatter
	IsStringPiece()
}

type StringSegment lex.Token

func (s *StringSegment) IsStringPiece() {}
func (s *StringSegment) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "StringSegment[%v]", s.Text)
}

// $NAME
// ${NAME}
// ${NAME:-WORD}
type ParameterExpansion struct {
	VarName  *lex.Token
	Operator *lex.Token
	Word     *lex.Token
}

func (p *ParameterExpansion) IsStringPiece() {}
func (p *ParameterExpansion) IsExpr()        {}
func (p *ParameterExpansion) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "ParameterExpansion[")
	if p.VarName != nil {
		fmt.Fprintf(f, "%v", p.VarName.Text)
	}
	if p.Operator != nil {
		fmt.Fprintf(f, "%v", p.Operator.Text)
	}
	if p.Word != nil {
		fmt.Fprintf(f, "%v", p.Word.Text)
	}
	fmt.Fprintf(f, "]")
}

func (p *ParameterExpansion) Parse(parser *Parser) error {
	if _, err := parser.ConsumeToken(lex.Dollar, nil); err != nil {
		return err
	}

	// $ { NAME OP WORD }
	if parser.Lexer.HasAnyToken(lex.LeftBrace) {
		parser.ConsumeToken(lex.LeftBrace, nil)
		if _, err := parser.ConsumeToken(lex.Name, &p.VarName); err != nil {
			return err
		}

		// TODO: implement the "OP WORD" bit of this

		if _, err := parser.ConsumeToken(lex.RightBrace, nil); err != nil {
			return err
		}
		return nil
	}

	// $ NAME
	if _, err := parser.ConsumeToken(lex.Name, &p.VarName); err != nil {
		return err
	}
	return nil
}
