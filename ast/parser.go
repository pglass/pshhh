package ast

import (
	"fmt"
	"io"
	//	"log"
	"psh/lex"
)

type Parser struct {
	Root  *GenericNode
	Lexer *lex.Lexer
	Debug bool
}

func NewParser(lexer *lex.Lexer) *Parser {
	return &Parser{
		Root:  &GenericNode{},
		Lexer: lexer,
		Debug: false,
	}
}

func (p *Parser) Parse() (Node, error) {
	// ignore leading spaces
	p.ConsumeWhile(lex.Space, lex.Newline)

	for {
		node, err := p.ParseNext()
		if err != nil && err != io.EOF {
			return nil, err
		} else if node != nil {
			p.Root.Children = append(p.Root.Children, node)
		} else {
			break
		}
	}
	return p.Root, nil
}

/* Parse the next command or expression. Returns (nil, nil) if nothing else
 * can be parsed (on EOF, for example). Returns an error if tokens remain in
 * the lexer but cannot be consumed */
func (p *Parser) ParseNext() (Node, error) {
	token, err := p.Lexer.Peek()
	if err != nil && err != io.EOF {
		return nil, err
	} else if err == io.EOF || token == nil {
		return nil, nil
	}

	var node Node = nil

	switch token.Type {
	case lex.AndIf, lex.OrIf, lex.DoubleQuote, lex.StringSegment:
		node, err = p.ParseExpr(token)
	case lex.For, lex.If, lex.Case, lex.While, lex.Until, lex.Word, lex.Name, lex.Number:
		command_list := NewCommandList()
		err = command_list.Parse(p)
		node = command_list
	default:
		err = fmt.Errorf("Syntax error near %v", token)
	}

	return node, err
}

func (p *Parser) ParseExpr(peekToken *lex.Token) (Expr, error) {
	var expr Expr = nil
	switch peekToken.Type {
	case lex.AndIf, lex.OrIf:
		if len(p.Root.Children) == 0 {
			return nil, fmt.Errorf(`Syntax error near %q`, peekToken.Text)
		}
		left := p.Root.Children[0]
		p.Root.Children = p.Root.Children[1:]
		expr = NewAndOrClause(left)
	case lex.DoubleQuote, lex.StringSegment:
		expr = NewString()
	default:
		return nil, nil
	}

	if err := expr.Parse(p); err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) ConsumeWhile(ttypes ...lex.TokenType) []*lex.Token {
	result := []*lex.Token{}
	for {
		token, err := p.Lexer.Peek()
		if err != nil || token == nil {
			break
		}

		found := false
		for _, ttype := range ttypes {
			if token.Type == ttype {
				found = true
				p.Lexer.Next()
				result = append(result, token)
				break
			}
		}

		if !found {
			break
		}
	}
	return result

}

func (p *Parser) ConsumeToken(ttype lex.TokenType, dst **lex.Token) (*lex.Token, error) {
	if tok, err := p.Lexer.Peek(); err != nil {
		return nil, err
	} else if tok == nil {
		return nil, nil
	} else if tok.Type != ttype {
		return nil, fmt.Errorf("Expected token of type %v (got %v)", ttype, tok)
	} else {
		p.Lexer.Next()
		if dst != nil {
			*dst = tok
		}
		return tok, err
	}
}

func (p *Parser) ConsumeAny(ttypes ...lex.TokenType) (*lex.Token, error) {
	if tok, err := p.Lexer.Peek(); err != nil {
		return nil, err
	} else if tok == nil {
		return nil, fmt.Errorf("Saw EOF but expected token of any type %v", ttypes)
	} else {
		for _, ttype := range ttypes {
			if tok.Type == ttype {
				p.Lexer.Next()
				return tok, nil
			}
		}
		return nil, fmt.Errorf("Unexpected token of type %v", tok.Type)
	}
}

func (p *Parser) ConsumeWordlist() []*lex.Token {
	tokens := []*lex.Token{}
	for _, tok := range p.ConsumeWhile(lex.Word, lex.Name, lex.Number, lex.Space) {
		if tok.Type != lex.Space {
			tokens = append(tokens, tok)
		}
	}
	return tokens
}
