package ast

import (
	"fmt"
	"psh/lex"
)

type SimpleCommand struct {
	Redirects []*IoRedirect
	// TODO: how do we store strings here?
	Words []*String
}

func NewSimpleCommand() *SimpleCommand {
	return &SimpleCommand{
		Redirects: []*IoRedirect{},
		Words:     []*String{},
	}
}

func (s *SimpleCommand) IsCommand() {}

func (s *SimpleCommand) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "SimpleCommand[")
	if len(s.Words) > 0 {
		fmt.Fprintf(f, "%v", s.Words[0])
		for _, word := range s.Words[1:] {
			fmt.Fprintf(f, " %v", word)
		}
	}

	if len(s.Redirects) > 0 {
		fmt.Fprintf(f, "%v", s.Redirects[0])
		for _, redirect := range s.Redirects {
			fmt.Fprintf(f, " %v", redirect)
		}
	}
	fmt.Fprintf(f, "]")
}

func (s *SimpleCommand) Parse(parser *Parser) error {
	// redirects can come anywhere within a simple command. That is, stuff like
	// the following is perfectly valid:
	//
	//    python >wumbo.txt -c "import sys; print sys.stdin.read()" <poo.txt
	//
	// We can basically just collect all the redirects into one list and
	// collect all the words into another list, as we see them (probably not
	// perfect, but it will do).
	for {
		parser.ConsumeWhile(lex.Space)

		if parsed, err := s.parseIoRedirect(parser); err != nil {
			return err
		} else if parsed {
			continue
		}

		parser.ConsumeWhile(lex.Space)

		if parser.Lexer.HasAnyToken(lex.DoubleQuote) {
			ast_str := NewString()
			if err := ast_str.Parse(parser); err != nil {
				return nil
			} else {
				s.Words = append(s.Words, ast_str)
				continue
			}
		}

		parser.ConsumeWhile(lex.Space)

		words := parser.ConsumeWordlist()
		for _, word := range words {
			// not the best
			ast_str := &String{Pieces: []StringPiece{(*StringSegment)(word)}}
			s.Words = append(s.Words, ast_str)
		}
		if len(words) > 0 {
			continue
		}

		// consumed no redirects or words
		break
	}

	return nil
}

func (s *SimpleCommand) parseIoRedirect(parser *Parser) (bool, error) {
	io_redirect := NewIoRedirect(parser)
	if err := io_redirect.Parse(); err != nil {
		return false, err
	} else if io_redirect.FilenameOrHereEnd != nil {
		s.Redirects = append(s.Redirects, io_redirect)
		return true, nil
	}
	return false, nil
}
