package ast

import (
	"fmt"
	"psh/lex"
)

type String struct {
	Pieces []StringPiece
}

func NewString() *String {
	return &String{Pieces: []StringPiece{}}
}

func NewStringFromTok(tok *lex.Token) *String {
	return &String{Pieces: []StringPiece{(*StringSegment)(tok)}}
}

func (s *String) IsExpr() {}

func (s *String) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "String[")
	for _, piece := range s.Pieces {
		fmt.Fprintf(f, "%v", piece)
	}
	fmt.Fprintf(f, "]")
}

func (s *String) Parse(parser *Parser) error {
	// a single-quoted string produces a single string segment
	if parser.Lexer.HasAnyToken(lex.StringSegment) {
		tok, _ := parser.Lexer.Next()
		piece := (*StringSegment)(tok)
		s.Pieces = append(s.Pieces, piece)
		return nil
	}

	// otherwise, parse a double-quoted string which is more involved
	if _, err := parser.ConsumeToken(lex.DoubleQuote, nil); err != nil {
		return err // if this
	}

	for {
		tok, err := parser.Lexer.Peek()
		if err != nil {
			return err
		}

		var piece StringPiece = nil
		if tok.Type == lex.StringSegment {
			parser.Lexer.Next()
			piece = (*StringSegment)(tok)
		} else if tok.Type == lex.Dollar {
			pe := &ParameterExpansion{}
			err := pe.Parse(parser)
			if err != nil {
				return err
			}
			piece = pe
		} else {
			break
		}

		s.Pieces = append(s.Pieces, piece)
	}

	if _, err := parser.ConsumeToken(lex.DoubleQuote, nil); err != nil {
		return err
	}
	return nil
}
