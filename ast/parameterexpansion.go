package ast

import (
	"fmt"

	"github.com/pglass/pshhh/lex"
)

type ParameterExpansion struct {
	VarName  *lex.Token
	Operator *lex.Token
	Word     *Str
}

func (p *ParameterExpansion) IsStrPiece() {}
func (p *ParameterExpansion) IsExpr()     {}

func (p *ParameterExpansion) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "ParameterExpansion[")
	if p.VarName != nil {
		fmt.Fprintf(f, "%v", p.VarName.Text)
	}
	if p.Operator != nil {
		fmt.Fprintf(f, "%v", p.Operator.Text)
	}
	if p.Word != nil {
		fmt.Fprintf(f, "%v", p.Word)
	}
	fmt.Fprintf(f, "]")
}

func (p *ParameterExpansion) Parse(parser *Parser) error {
	if _, err := parser.ConsumeToken(lex.Dollar, nil); err != nil {
		return err
	}

	if parser.Lexer.HasAnyToken(lex.LeftBrace) {
		// $ { NAME OP WORD }
		parser.ConsumeToken(lex.LeftBrace, nil)
		if _, err := parser.ConsumeToken(lex.Name, &p.VarName); err != nil {
			return err
		}

		if tok, _ := parser.ConsumeAny(lex.PARAMETER_EXPANSION_TTYPES...); tok != nil {
			p.Operator = tok

			str := NewStr()
			err := str.Parse(parser)
			if err != nil {
				return err
			}
			p.Word = str
		}

		if _, err := parser.ConsumeToken(lex.RightBrace, nil); err != nil {
			return err
		}
		return nil
	} else if _, err := parser.ConsumeToken(lex.Name, &p.VarName); err != nil {
		// $ NAME
		return err
	}
	return nil
}
