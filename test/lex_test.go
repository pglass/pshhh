package test

import (
	"testing"

	"github.com/pglass/pshhh/lex"
	"github.com/stretchr/testify/assert"
)

func lex_input(input string) []lex.Token {
	lexer := lex.NewLexer(input)

	tokens := []lex.Token{}
	for {
		tok := lexer.Next()
		tokens = append(tokens, tok)

		switch tok.Type {
		case lex.EOF, lex.ERROR:
			return tokens
		}
	}
	return tokens
}

type lexData struct {
	Input  string
	Tokens []lex.Token
}

func run_lex_test(t *testing.T, data lexData) {
	t.Logf("Input: %q", data.Input)
	t.Logf("Output (expected): %v", data.Tokens)

	tokens := lex_input(data.Input)

	t.Logf("Output (received): %v", tokens)

	assert.Equal(t, data.Tokens, tokens)
}

func TestLexer(t *testing.T) {
	for _, data := range LEX_CASES {
		t.Run(data.Input, func(t *testing.T) { run_lex_test(t, data) })
	}
}

var LEX_CASES = []lexData{
	lexData{
		Input: "abc foo.txt ./main --verbose=1",
		Tokens: []lex.Token{
			lex.Token{lex.Name, "abc", 0, 1},
			lex.Token{lex.Space, " ", 3, 1},
			lex.Token{lex.Name, "foo.txt", 4, 1},
			lex.Token{lex.Space, " ", 11, 1},
			lex.Token{lex.Name, "./main", 12, 1},
			lex.Token{lex.Space, " ", 18, 1},
			lex.Token{lex.Name, "--verbose=1", 19, 1},
			lex.Token{lex.EOF, "", 30, 1},
		},
	},
	lexData{
		Input: "1abc 234a 5b7",
		Tokens: []lex.Token{
			lex.Token{lex.Word, "1abc", 0, 1},
			lex.Token{lex.Space, " ", 4, 1},
			lex.Token{lex.Word, "234a", 5, 1},
			lex.Token{lex.Space, " ", 9, 1},
			lex.Token{lex.Word, "5b7", 10, 1},
			lex.Token{lex.EOF, "", 13, 1},
		},
	},
	lexData{
		Input: "123",
		Tokens: []lex.Token{
			lex.Token{lex.Number, "123", 0, 1},
			lex.Token{lex.EOF, "", 3, 1},
		},
	},
	lexData{
		Input: "a23 4ef ghi",
		Tokens: []lex.Token{
			lex.Token{lex.Name, "a23", 0, 1},
			lex.Token{lex.Space, " ", 3, 1},
			lex.Token{lex.Word, "4ef", 4, 1},
			lex.Token{lex.Space, " ", 7, 1},
			lex.Token{lex.Name, "ghi", 8, 1},
			lex.Token{lex.EOF, "", 11, 1},
		},
	},
	lexData{
		Input: "\n\n\n\n",
		Tokens: []lex.Token{
			lex.Token{lex.Newline, "\n", 0, 1},
			lex.Token{lex.Newline, "\n", 1, 2},
			lex.Token{lex.Newline, "\n", 2, 3},
			lex.Token{lex.Newline, "\n", 3, 4},
			lex.Token{lex.EOF, "", 4, 5},
		},
	},
	lexData{
		Input: "a \n 2b\nc\n ",
		Tokens: []lex.Token{
			lex.Token{lex.Name, "a", 0, 1},
			lex.Token{lex.Space, " ", 1, 1},
			lex.Token{lex.Newline, "\n", 2, 1},
			lex.Token{lex.Space, " ", 3, 2},
			lex.Token{lex.Word, "2b", 4, 2},
			lex.Token{lex.Newline, "\n", 6, 2},
			lex.Token{lex.Name, "c", 7, 3},
			lex.Token{lex.Newline, "\n", 8, 3},
			lex.Token{lex.Space, " ", 9, 4},
			lex.Token{lex.EOF, "", 10, 4},
		},
	},
	lexData{
		Input: "a  &&  b  ||  c",
		Tokens: []lex.Token{
			lex.Token{lex.Name, "a", 0, 1},
			lex.Token{lex.Space, "  ", 1, 1},
			lex.Token{lex.AndIf, "&&", 3, 1},
			lex.Token{lex.Space, "  ", 5, 1},
			lex.Token{lex.Name, "b", 7, 1},
			lex.Token{lex.Space, "  ", 8, 1},
			lex.Token{lex.OrIf, "||", 10, 1},
			lex.Token{lex.Space, "  ", 12, 1},
			lex.Token{lex.Name, "c", 14, 1},
			lex.Token{lex.EOF, "", 15, 1},
		},
	},
	lexData{
		Input: "for var in items; do echo; done",
		Tokens: []lex.Token{
			lex.Token{lex.For, "for", 0, 1},
			lex.Token{lex.Space, " ", 3, 1},
			lex.Token{lex.Name, "var", 4, 1},
			lex.Token{lex.Space, " ", 7, 1},
			lex.Token{lex.In, "in", 8, 1},
			lex.Token{lex.Space, " ", 10, 1},
			lex.Token{lex.Name, "items", 11, 1},
			lex.Token{lex.Semi, ";", 16, 1},
			lex.Token{lex.Space, " ", 17, 1},
			lex.Token{lex.Do, "do", 18, 1},
			lex.Token{lex.Space, " ", 20, 1},
			lex.Token{lex.Name, "echo", 21, 1},
			lex.Token{lex.Semi, ";", 25, 1},
			lex.Token{lex.Space, " ", 26, 1},
			lex.Token{lex.Done, "done", 27, 1},
			lex.Token{lex.EOF, "", 31, 1},
		},
	},
	lexData{
		Input: "&&&",
		Tokens: []lex.Token{
			lex.Token{lex.AndIf, "&&", 0, 1},
			lex.Token{lex.Ampersand, "&", 2, 1},
			lex.Token{lex.EOF, "", 3, 1},
		},
	},
	lexData{
		Input: "<<<<-",
		Tokens: []lex.Token{
			lex.Token{lex.DoubleLess, "<<", 0, 1},
			lex.Token{lex.DoubleLessDash, "<<-", 2, 1},
			lex.Token{lex.EOF, "", 5, 1},
		},
	},
	lexData{
		Input: "< >",
		Tokens: []lex.Token{
			lex.Token{lex.Less, "<", 0, 1},
			lex.Token{lex.Space, " ", 1, 1},
			lex.Token{lex.Great, ">", 2, 1},
			lex.Token{lex.EOF, "", 3, 1},
		},
	},
	lexData{
		Input: "<<<>>>",
		Tokens: []lex.Token{
			lex.Token{lex.DoubleLess, "<<", 0, 1},
			lex.Token{lex.LessGreat, "<>", 2, 1},
			lex.Token{lex.DoubleGreat, ">>", 4, 1},
			lex.Token{lex.EOF, "", 6, 1},
		},
	},
	lexData{
		Input: "<<>>",
		Tokens: []lex.Token{
			lex.Token{lex.DoubleLess, "<<", 0, 1},
			lex.Token{lex.DoubleGreat, ">>", 2, 1},
			lex.Token{lex.EOF, "", 4, 1},
		},
	},
	lexData{
		Input: "$",
		Tokens: []lex.Token{
			lex.Token{lex.Dollar, "$", 0, 1},
			lex.Token{lex.EOF, "", 1, 1},
		},
	},
	lexData{
		Input: "''",
		Tokens: []lex.Token{
			lex.Token{lex.SingleQuote, "'", 0, 1},
			lex.Token{lex.StringSegment, "", 1, 1},
			lex.Token{lex.SingleQuote, "'", 1, 1},
			lex.Token{lex.EOF, "", 2, 1},
		},
	},
	lexData{
		// single-quoted strings do not understand backslash escapes
		Input: `'\' '\'''`,
		Tokens: []lex.Token{
			lex.Token{lex.SingleQuote, "'", 0, 1},
			lex.Token{lex.StringSegment, `\`, 1, 1},
			lex.Token{lex.SingleQuote, "'", 2, 1},
			lex.Token{lex.Space, " ", 3, 1},
			lex.Token{lex.SingleQuote, "'", 4, 1},
			lex.Token{lex.StringSegment, `\`, 5, 1},
			lex.Token{lex.SingleQuote, "'", 6, 1},
			lex.Token{lex.SingleQuote, "'", 7, 1},
			lex.Token{lex.StringSegment, "", 8, 1},
			lex.Token{lex.SingleQuote, "'", 8, 1},
			lex.Token{lex.EOF, "", 9, 1},
		},
	},
	lexData{
		// single-quoted strings do not understand parameter expansions
		Input: `'$WUMBO'`,
		Tokens: []lex.Token{
			lex.Token{lex.SingleQuote, "'", 0, 1},
			lex.Token{lex.StringSegment, "$WUMBO", 1, 1},
			lex.Token{lex.SingleQuote, "'", 7, 1},
			lex.Token{lex.EOF, "", 8, 1},
		},
	},
	lexData{
		Input: `"$"`,
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.Dollar, "$", 1, 1},
			lex.Token{lex.DoubleQuote, `"`, 2, 1},
			lex.Token{lex.EOF, "", 3, 1},
		},
	},
	lexData{
		Input: `"$WUMBO"`,
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.Dollar, "$", 1, 1},
			lex.Token{lex.Name, "WUMBO", 2, 1},
			lex.Token{lex.DoubleQuote, `"`, 7, 1},
			lex.Token{lex.EOF, "", 8, 1},
		},
	},
	lexData{
		Input: "\"$WUMBO \n $MINI\"",
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.Dollar, "$", 1, 1},
			lex.Token{lex.Name, "WUMBO", 2, 1},
			lex.Token{lex.StringSegment, " \n ", 7, 1},
			lex.Token{lex.Dollar, "$", 10, 2},
			lex.Token{lex.Name, "MINI", 11, 2},
			lex.Token{lex.DoubleQuote, `"`, 15, 2},
			lex.Token{lex.EOF, "", 16, 2},
		},
	},
	lexData{
		Input: `"\$"`,
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.StringSegment, "$", 1, 1},
			lex.Token{lex.DoubleQuote, `"`, 3, 1},
			lex.Token{lex.EOF, "", 4, 1},
		},
	},
	lexData{
		Input: `"\\"`,
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.StringSegment, `\`, 1, 1},
			lex.Token{lex.DoubleQuote, `"`, 3, 1},
			lex.Token{lex.EOF, "", 4, 1},
		},
	},
	lexData{
		// "\`" -- checking an escaped backtick within a double-quoted string
		Input: "\"\\`\"",
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.StringSegment, "`", 1, 1},
			lex.Token{lex.DoubleQuote, `"`, 3, 1},
			lex.Token{lex.EOF, "", 4, 1},
		},
	},
	lexData{
		Input: `"\""`,
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.StringSegment, `"`, 1, 1},
			lex.Token{lex.DoubleQuote, `"`, 3, 1},
			lex.Token{lex.EOF, "", 4, 1},
		},
	},
	lexData{
		Input: `"$\$$\$"`,
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.Dollar, `$`, 1, 1},
			lex.Token{lex.StringSegment, `$`, 2, 1},
			lex.Token{lex.Dollar, `$`, 4, 1},
			lex.Token{lex.StringSegment, `$`, 5, 1},
			lex.Token{lex.DoubleQuote, `"`, 7, 1},
			lex.Token{lex.EOF, "", 8, 1},
		},
	},
	lexData{
		Input: `"unclosed-string`,
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.ERROR, "Unclosed string", 1, 1},
		},
	},
	lexData{
		Input: `"${WUMBO}"`,
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.Dollar, "$", 1, 1},
			lex.Token{lex.LeftBrace, "{", 2, 1},
			lex.Token{lex.Name, "WUMBO", 3, 1},
			lex.Token{lex.RightBrace, "}", 8, 1},
			lex.Token{lex.DoubleQuote, `"`, 9, 1},
			lex.Token{lex.EOF, "", 10, 1},
		},
	},
	lexData{
		Input: `${MYKEY:-myval}`,
		Tokens: []lex.Token{
			lex.Token{lex.Dollar, "$", 0, 1},
			lex.Token{lex.LeftBrace, "{", 1, 1},
			lex.Token{lex.Name, "MYKEY", 2, 1},
			lex.Token{lex.ColonDash, ":-", 7, 1},
			lex.Token{lex.Name, "myval", 9, 1},
			lex.Token{lex.RightBrace, "}", 14, 1},
			lex.Token{lex.EOF, "", 15, 1},
		},
	},
	lexData{
		Input: `${MYKEY:="myval"}`,
		Tokens: []lex.Token{
			lex.Token{lex.Dollar, "$", 0, 1},
			lex.Token{lex.LeftBrace, "{", 1, 1},
			lex.Token{lex.Name, "MYKEY", 2, 1},
			lex.Token{lex.ColonEquals, ":=", 7, 1},
			lex.Token{lex.DoubleQuote, `"`, 9, 1},
			lex.Token{lex.StringSegment, "myval", 10, 1},
			lex.Token{lex.DoubleQuote, `"`, 15, 1},
			lex.Token{lex.RightBrace, "}", 16, 1},
			lex.Token{lex.EOF, "", 17, 1},
		},
	},
	lexData{
		Input: `${MYKEY:?"myval"}`,
		Tokens: []lex.Token{
			lex.Token{lex.Dollar, "$", 0, 1},
			lex.Token{lex.LeftBrace, "{", 1, 1},
			lex.Token{lex.Name, "MYKEY", 2, 1},
			lex.Token{lex.ColonQuestion, ":?", 7, 1},
			lex.Token{lex.DoubleQuote, `"`, 9, 1},
			lex.Token{lex.StringSegment, "myval", 10, 1},
			lex.Token{lex.DoubleQuote, `"`, 15, 1},
			lex.Token{lex.RightBrace, "}", 16, 1},
			lex.Token{lex.EOF, "", 17, 1},
		},
	},
	lexData{
		Input: `"${MYKEY:+myval}"`,
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.Dollar, "$", 1, 1},
			lex.Token{lex.LeftBrace, "{", 2, 1},
			lex.Token{lex.Name, "MYKEY", 3, 1},
			lex.Token{lex.ColonPlus, ":+", 8, 1},
			lex.Token{lex.Name, "myval", 10, 1},
			lex.Token{lex.RightBrace, "}", 15, 1},
			lex.Token{lex.DoubleQuote, `"`, 16, 1},
			lex.Token{lex.EOF, "", 17, 1},
		},
	},
	lexData{
		Input: `"${MYKEY+"myval"}"`,
		Tokens: []lex.Token{
			lex.Token{lex.DoubleQuote, `"`, 0, 1},
			lex.Token{lex.Dollar, "$", 1, 1},
			lex.Token{lex.LeftBrace, "{", 2, 1},
			lex.Token{lex.Name, "MYKEY", 3, 1},
			lex.Token{lex.Plus, "+", 8, 1},
			lex.Token{lex.DoubleQuote, `"`, 9, 1},
			lex.Token{lex.StringSegment, "myval", 10, 1},
			lex.Token{lex.DoubleQuote, `"`, 15, 1},
			lex.Token{lex.RightBrace, "}", 16, 1},
			lex.Token{lex.DoubleQuote, `"`, 17, 1},
			lex.Token{lex.EOF, "", 18, 1},
		},
	},
	lexData{
		Input: `echo ${MYKEY-myval}`,
		Tokens: []lex.Token{
			lex.Token{lex.Name, "echo", 0, 1},
			lex.Token{lex.Space, " ", 4, 1},
			lex.Token{lex.Dollar, "$", 5, 1},
			lex.Token{lex.LeftBrace, "{", 6, 1},
			lex.Token{lex.Name, "MYKEY", 7, 1},
			lex.Token{lex.Dash, "-", 12, 1},
			lex.Token{lex.Name, "myval", 13, 1},
			lex.Token{lex.RightBrace, "}", 18, 1},
			lex.Token{lex.EOF, "", 19, 1},
		},
	},
}
