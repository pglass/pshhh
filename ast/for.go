package ast

import (
	"fmt"
	"psh/lex"
)

type ForClause struct {
	LoopVar *lex.Token
	In      *lex.Token

	// an OPTIONAL list of Words and Names
	// note: expansions must occur before parsing. If you have
	//
	//		for x in `seq 1 3`
	//
	// The expression `seq 1 3` is command substitution (in bash). In order
	// to execute the for loop, we would first tokenize. We would go through
	// the tokens and run the command `seq 1 3` to get "1\n2\n3\n", which we
	// would tokenize to tokens "1", "\n", "2", "\n", "3", "\n". These tokens
	// would then replace `seq 1 3` in our token stream.
	Wordlist []lex.Token

	DoClause *DoClause
}

func NewForClause() *ForClause {
	return &ForClause{
		Wordlist: []lex.Token{},
		DoClause: NewDoClause(),
	}
}

func (f *ForClause) IsCommand() {}

func (f *ForClause) Format(fs fmt.State, c rune) {
	fmt.Fprintf(fs, "For %q", f.LoopVar.Text)
	if f.In != nil {
		fmt.Fprintf(fs, " in %v", f.Wordlist)
	}
}

func (f *ForClause) Parse(parser *Parser) error {
	// "for"
	if _, err := parser.ConsumeToken(lex.For, nil); err != nil {
		return err
	}

	parser.ConsumeWhile(lex.Space)

	// "<var>"
	if _, err := parser.ConsumeToken(lex.Name, &f.LoopVar); err != nil {
		return err
	}

	parser.ConsumeWhile(lex.Space)

	// if there is no in-clause, err is nil. only error on real errors.
	if err := f.parseOptionalInClause(parser); err != nil {
		return err
	}

	// "do ... done"
	if err := f.DoClause.Parse(parser); err != nil {
		return err
	}

	return nil
}

func (f *ForClause) parseOptionalInClause(parser *Parser) error {
	if !parser.Lexer.HasAnyToken(lex.In) {
		return nil
	}

	// consume "in"
	if _, err := parser.ConsumeToken(lex.In, &f.In); err != nil {
		return err
	}

	// " "
	parser.ConsumeWhile(lex.Space)

	// consume and store the word list
	f.Wordlist = append(f.Wordlist, parser.ConsumeWordlist()...)
	// for _, tok := range parser.ConsumeWhile(lex.Word, lex.Name, lex.Space) {
	// 	switch tok.Type {
	// 	case lex.Word, lex.Name:
	// 		f.Wordlist = append(f.Wordlist, tok)
	// 	}
	// }

	// " "
	parser.ConsumeWhile(lex.Space)

	// the separator is either ';' or '\n'
	// if there is no in clause, we don't *need* a separator
	//
	//   # these are all OKAY
	//   for x do echo $x; done
	//	 for x do
	//		echo $x; done
	//	 for x
	//	 do echo $x; done
	//   for x; do echo $x; done
	//   for x;
	//	 do echo $x; done
	//
	// if there is an "in" clause, we must have either a separator (; or \n)
	//
	//	 # these are OKAY
	//	 for x in; do echo done
	//	 for x in 1 2 3; do echo done
	//
	//	 # these are NOT okay
	//	 for x in do echo $x; done
	//   for x in 1 2 3 do echo $x; done
	//
	has_separator := parser.Lexer.HasAnyToken(lex.Semi, lex.Newline)
	if f.In != nil && !has_separator {
		return fmt.Errorf(`Separator (';' or newline) is required in for loop with "in" clause`)
	}
	if has_separator {
		parser.Lexer.Next()
	}

	parser.ConsumeWhile(lex.Space, lex.Newline)

	return nil
}
