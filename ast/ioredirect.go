package ast

import (
	"fmt"
	"psh/lex"
)

type IoRedirect struct {
	parser            *Parser
	IoNumber          *lex.Token
	IoOperator        *lex.Token
	FilenameOrHereEnd *lex.Token
}

func NewIoRedirect(parser *Parser) *IoRedirect {
	return &IoRedirect{
		parser:            parser,
		IoNumber:          nil,
		IoOperator:        nil,
		FilenameOrHereEnd: nil,
	}
}

func (i *IoRedirect) Format(f fmt.State, c rune) {
	if i.IoNumber != nil {
		fmt.Fprintf(f, "%v", i.IoNumber.Text)
	}
	fmt.Fprintf(f, "%v", i.IoOperator.Text)
	fmt.Fprintf(f, "%v", i.FilenameOrHereEnd.Text)
}

func (i *IoRedirect) Parse() error {
	// consume the optional IO_NUMBER
	// note: NO SPACES between the IO number and the operator
	i.parser.ConsumeToken(lex.Number, &i.IoNumber)

	// consume the required operator
	tok, _ := i.parser.ConsumeAny(lex.Less, lex.LessAnd, lex.Great,
		lex.GreatAnd, lex.DoubleGreat, lex.LessGreat, lex.Clobber,
		lex.DoubleLess, lex.DoubleLessDash)
	i.IoOperator = tok

	// if there is no redirect operator, unconsume the io number.
	// (this is a "tried, failed, backtrack" situation)
	if i.IoOperator == nil {
		if i.IoNumber != nil {
			i.parser.Lexer.Unread(*i.IoNumber)
		}
		return nil
	}

	i.parser.ConsumeWhile(lex.Space)

	// at this point, we've seen the operator, so there must be a word following
	tok, err := i.parser.ConsumeAny(lex.Word, lex.Name, lex.Number)
	if err != nil {
		return err
	}
	i.FilenameOrHereEnd = tok

	return nil
}
