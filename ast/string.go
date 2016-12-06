package ast

import (
	"fmt"

	"github.com/pglass/pshhh/lex"
)

// a Str is composed of a sequence of pieces. each piece can be:
//
//   1. RawStr (raw text)
//   2. ParameterExpansion (a substitution)
//
type Str struct {
	Pieces []StrPiece
}

func NewStr() *Str {
	return &Str{Pieces: []StrPiece{}}
}

func NewStrFromTok(tok lex.Token) *Str {
	return &Str{Pieces: []StrPiece{RawStr(tok.Text)}}
}

func (s *Str) IsExpr() {}

func (s *Str) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "Str[")
	for _, piece := range s.Pieces {
		fmt.Fprintf(f, "%v", piece)
	}
	fmt.Fprintf(f, "]")
}

func (s *Str) Parse(parser *Parser) error {
	// a single-quoted string produces a single string segment
	tok := parser.Lexer.Peek()
	switch tok.Type {
	case lex.SingleQuote:
		return s.parseSingleQuotedString(parser)
	case lex.DoubleQuote:
		return s.parseDoubleQuotedString(parser)
	case lex.Dollar:
		return s.parseParamExpansion(parser)
	case lex.Word, lex.Name, lex.Number:
		parser.Lexer.Next()
		s.Pieces = append(s.Pieces, RawStr(tok.Text))
		return nil
	default:
		return fmt.Errorf("Failed to parse a Str [bug?]")
	}
}

func (s *Str) parseSingleQuotedString(parser *Parser) error {
	if parser.Lexer.Next().Type != lex.SingleQuote {
		return fmt.Errorf("Expected single quote [bug!]")
	}

	// this is optional. we could have an empty string.
	if parser.Lexer.HasAnyToken(lex.StringSegment) {
		tok := parser.Lexer.Next()
		piece := RawStr(tok.Text)
		s.Pieces = append(s.Pieces, piece)
		return nil
	}

	if parser.Lexer.HasAnyToken(lex.SingleQuote) {
		parser.Lexer.Next()
	} else {
		return fmt.Errorf("Unclosed string (expected ')")
	}
	return nil
}

func (s *Str) parseParamExpansion(parser *Parser) error {
	pe := &ParameterExpansion{}
	if err := pe.Parse(parser); err != nil {
		return err
	}
	s.Pieces = append(s.Pieces, StrPiece(pe))
	return nil
}

func (s *Str) parseDoubleQuotedString(parser *Parser) error {
	if _, err := parser.ConsumeToken(lex.DoubleQuote, nil); err != nil {
		return err
	}

	for {
		tok := parser.Lexer.Peek()
		switch tok.Type {
		case lex.EOF:
			return fmt.Errorf("Unclosed string (expected \")")
		case lex.ERROR:
			return fmt.Errorf("%v", tok)
		case lex.DoubleQuote:
			parser.Lexer.Next()
			return nil
		case lex.StringSegment:
			parser.Lexer.Next()
			s.Pieces = append(s.Pieces, RawStr(tok.Text))
		case lex.Dollar:
			if err := s.parseParamExpansion(parser); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unexpected token in string: %v", tok)
		}
	}

	return nil

}
