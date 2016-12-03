package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"psh/lex"
	"strings"
	"testing"
)

func lex_input(input string) ([]*lex.Token, error) {
	reader := strings.NewReader(input)
	lexer := lex.NewLexer(reader)

	tokens := []*lex.Token{}
	var token *lex.Token = nil
	var err error = nil
	for {
		token, err = lexer.Next()
		if token != nil {
			tokens = append(tokens, token)
		} else if err != nil || token == nil {
			break
		}
	}
	return tokens, err
}

type lexData struct {
	Input  string
	Tokens []*lex.Token
	Error  error
}

func run_lex_test(t *testing.T, data lexData) {
	t.Logf("Input: %q", data.Input)
	t.Logf("Error (expected): %v", data.Error)
	t.Logf("Output (expected): %v", data.Tokens)

	tokens, err := lex_input(data.Input)

	t.Logf("Error (recieved): %v", err)
	t.Logf("Output (received): %v", tokens)

	if data.Error != nil {
		assert.Equal(t, data.Error, err)
	} else {
		assert.Nil(t, err)
	}

	if data.Tokens != nil {
		assert.Equal(t, data.Tokens, tokens)
	}
}

func TestLexer(t *testing.T) {
	for _, data := range LEX_CASES {
		t.Run(data.Input, func(t *testing.T) { run_lex_test(t, data) })
	}
}

var LEX_CASES = []lexData{
	lexData{
		Input: "abc foo.txt ./main",
		Tokens: []*lex.Token{
			&lex.Token{lex.Name, "abc", lex.Location{0, 0}, lex.Location{0, 3}},
			&lex.Token{lex.Space, " ", lex.Location{0, 3}, lex.Location{0, 4}},
			&lex.Token{lex.Name, "foo.txt", lex.Location{0, 4}, lex.Location{0, 11}},
			&lex.Token{lex.Space, " ", lex.Location{0, 11}, lex.Location{0, 12}},
			&lex.Token{lex.Name, "./main", lex.Location{0, 12}, lex.Location{0, 18}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "1abc 234a 5b7",
		Tokens: []*lex.Token{
			&lex.Token{lex.Word, "1abc", lex.Location{0, 0}, lex.Location{0, 4}},
			&lex.Token{lex.Space, " ", lex.Location{0, 4}, lex.Location{0, 5}},
			&lex.Token{lex.Word, "234a", lex.Location{0, 5}, lex.Location{0, 9}},
			&lex.Token{lex.Space, " ", lex.Location{0, 9}, lex.Location{0, 10}},
			&lex.Token{lex.Word, "5b7", lex.Location{0, 10}, lex.Location{0, 13}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "123",
		Tokens: []*lex.Token{
			&lex.Token{lex.Number, "123", lex.Location{0, 0}, lex.Location{0, 3}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "a23 4ef ghi",
		Tokens: []*lex.Token{
			&lex.Token{lex.Name, "a23", lex.Location{0, 0}, lex.Location{0, 3}},
			&lex.Token{lex.Space, " ", lex.Location{0, 3}, lex.Location{0, 4}},
			&lex.Token{lex.Word, "4ef", lex.Location{0, 4}, lex.Location{0, 7}},
			&lex.Token{lex.Space, " ", lex.Location{0, 7}, lex.Location{0, 8}},
			&lex.Token{lex.Name, "ghi", lex.Location{0, 8}, lex.Location{0, 11}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "\n\n\n\n",
		Tokens: []*lex.Token{
			&lex.Token{lex.Newline, "\n", lex.Location{0, 0}, lex.Location{1, 0}},
			&lex.Token{lex.Newline, "\n", lex.Location{1, 0}, lex.Location{2, 0}},
			&lex.Token{lex.Newline, "\n", lex.Location{2, 0}, lex.Location{3, 0}},
			&lex.Token{lex.Newline, "\n", lex.Location{3, 0}, lex.Location{4, 0}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "a \n 2b\nc\n ",
		Tokens: []*lex.Token{
			&lex.Token{lex.Name, "a", lex.Location{0, 0}, lex.Location{0, 1}},
			&lex.Token{lex.Space, " ", lex.Location{0, 1}, lex.Location{0, 2}},
			&lex.Token{lex.Newline, "\n", lex.Location{0, 2}, lex.Location{1, 0}},
			&lex.Token{lex.Space, " ", lex.Location{1, 0}, lex.Location{1, 1}},
			&lex.Token{lex.Word, "2b", lex.Location{1, 1}, lex.Location{1, 3}},
			&lex.Token{lex.Newline, "\n", lex.Location{1, 3}, lex.Location{2, 0}},
			&lex.Token{lex.Name, "c", lex.Location{2, 0}, lex.Location{2, 1}},
			&lex.Token{lex.Newline, "\n", lex.Location{2, 1}, lex.Location{3, 0}},
			&lex.Token{lex.Space, " ", lex.Location{3, 0}, lex.Location{3, 1}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "a  &&  b  ||  c",
		Tokens: []*lex.Token{
			&lex.Token{lex.Name, "a", lex.Location{0, 0}, lex.Location{0, 1}},
			&lex.Token{lex.Space, "  ", lex.Location{0, 1}, lex.Location{0, 3}},
			&lex.Token{lex.AndIf, "&&", lex.Location{0, 3}, lex.Location{0, 5}},
			&lex.Token{lex.Space, "  ", lex.Location{0, 5}, lex.Location{0, 7}},
			&lex.Token{lex.Name, "b", lex.Location{0, 7}, lex.Location{0, 8}},
			&lex.Token{lex.Space, "  ", lex.Location{0, 8}, lex.Location{0, 10}},
			&lex.Token{lex.OrIf, "||", lex.Location{0, 10}, lex.Location{0, 12}},
			&lex.Token{lex.Space, "  ", lex.Location{0, 12}, lex.Location{0, 14}},
			&lex.Token{lex.Name, "c", lex.Location{0, 14}, lex.Location{0, 15}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "for var in items; do echo; done",
		Tokens: []*lex.Token{
			&lex.Token{lex.For, "for", lex.Location{0, 0}, lex.Location{0, 3}},
			&lex.Token{lex.Space, " ", lex.Location{0, 3}, lex.Location{0, 4}},
			&lex.Token{lex.Name, "var", lex.Location{0, 4}, lex.Location{0, 7}},
			&lex.Token{lex.Space, " ", lex.Location{0, 7}, lex.Location{0, 8}},
			&lex.Token{lex.In, "in", lex.Location{0, 8}, lex.Location{0, 10}},
			&lex.Token{lex.Space, " ", lex.Location{0, 10}, lex.Location{0, 11}},
			&lex.Token{lex.Name, "items", lex.Location{0, 11}, lex.Location{0, 16}},
			&lex.Token{lex.Semi, ";", lex.Location{0, 16}, lex.Location{0, 17}},
			&lex.Token{lex.Space, " ", lex.Location{0, 17}, lex.Location{0, 18}},
			&lex.Token{lex.Do, "do", lex.Location{0, 18}, lex.Location{0, 20}},
			&lex.Token{lex.Space, " ", lex.Location{0, 20}, lex.Location{0, 21}},
			&lex.Token{lex.Name, "echo", lex.Location{0, 21}, lex.Location{0, 25}},
			&lex.Token{lex.Semi, ";", lex.Location{0, 25}, lex.Location{0, 26}},
			&lex.Token{lex.Space, " ", lex.Location{0, 26}, lex.Location{0, 27}},
			&lex.Token{lex.Done, "done", lex.Location{0, 27}, lex.Location{0, 31}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "&&&",
		Tokens: []*lex.Token{
			&lex.Token{lex.AndIf, "&&", lex.Location{0, 0}, lex.Location{0, 2}},
			&lex.Token{lex.Ampersand, "&", lex.Location{0, 2}, lex.Location{0, 3}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "<<<<-",
		Tokens: []*lex.Token{
			&lex.Token{lex.DoubleLess, "<<", lex.Location{0, 0}, lex.Location{0, 2}},
			&lex.Token{lex.DoubleLessDash, "<<-", lex.Location{0, 2}, lex.Location{0, 5}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "<-",
		Tokens: []*lex.Token{
			&lex.Token{lex.Less, "<", lex.Location{0, 0}, lex.Location{0, 1}},
		},
		Error: fmt.Errorf(`Syntax error near '-'`),
	},
	lexData{
		Input: "< >",
		Tokens: []*lex.Token{
			&lex.Token{lex.Less, "<", lex.Location{0, 0}, lex.Location{0, 1}},
			&lex.Token{lex.Space, " ", lex.Location{0, 1}, lex.Location{0, 2}},
			&lex.Token{lex.Great, ">", lex.Location{0, 2}, lex.Location{0, 3}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "<>",
		Tokens: []*lex.Token{
			&lex.Token{lex.LessGreat, "<>", lex.Location{0, 0}, lex.Location{0, 2}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "<<>>",
		Tokens: []*lex.Token{
			&lex.Token{lex.DoubleLess, "<<", lex.Location{0, 0}, lex.Location{0, 2}},
			&lex.Token{lex.DoubleGreat, ">>", lex.Location{0, 2}, lex.Location{0, 4}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "$",
		Tokens: []*lex.Token{
			&lex.Token{lex.Dollar, "$", lex.Location{0, 0}, lex.Location{0, 1}},
		},
		Error: io.EOF,
	},
	lexData{
		// single-quoted strings do not understand backslash escapes
		Input: `'' '\' '\'''`,
		Tokens: []*lex.Token{
			&lex.Token{lex.StringSegment, "", lex.Location{0, 0}, lex.Location{0, 2}},
			&lex.Token{lex.Space, " ", lex.Location{0, 2}, lex.Location{0, 3}},
			&lex.Token{lex.StringSegment, `\`, lex.Location{0, 3}, lex.Location{0, 6}},
			&lex.Token{lex.Space, " ", lex.Location{0, 6}, lex.Location{0, 7}},
			&lex.Token{lex.StringSegment, `\`, lex.Location{0, 7}, lex.Location{0, 10}},
			&lex.Token{lex.StringSegment, "", lex.Location{0, 10}, lex.Location{0, 12}},
		},
		Error: io.EOF,
	},
	lexData{
		// single-quoted strings do not understand parameter expansions
		Input: `'$WUMBO'`,
		Tokens: []*lex.Token{
			&lex.Token{lex.StringSegment, `$WUMBO`, lex.Location{0, 0}, lex.Location{0, 8}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: `"$WUMBO"`,
		Tokens: []*lex.Token{
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 0}, lex.Location{0, 1}},
			&lex.Token{lex.Dollar, `$`, lex.Location{0, 1}, lex.Location{0, 2}},
			&lex.Token{lex.Name, `WUMBO`, lex.Location{0, 2}, lex.Location{0, 7}},
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 7}, lex.Location{0, 8}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: "\"$WUMBO \n $MINI\"",
		Tokens: []*lex.Token{
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 0}, lex.Location{0, 1}},
			&lex.Token{lex.Dollar, `$`, lex.Location{0, 1}, lex.Location{0, 2}},
			&lex.Token{lex.Name, `WUMBO`, lex.Location{0, 2}, lex.Location{0, 7}},
			&lex.Token{lex.StringSegment, " \n ", lex.Location{0, 7}, lex.Location{1, 1}},
			&lex.Token{lex.Dollar, `$`, lex.Location{1, 1}, lex.Location{1, 2}},
			&lex.Token{lex.Name, `MINI`, lex.Location{1, 2}, lex.Location{1, 6}},
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{1, 6}, lex.Location{1, 7}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: `"\$"`,
		Tokens: []*lex.Token{
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 0}, lex.Location{0, 1}},
			&lex.Token{lex.StringSegment, `$`, lex.Location{0, 1}, lex.Location{0, 3}},
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 3}, lex.Location{0, 4}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: `"\\"`,
		Tokens: []*lex.Token{
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 0}, lex.Location{0, 1}},
			&lex.Token{lex.StringSegment, `\`, lex.Location{0, 1}, lex.Location{0, 3}},
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 3}, lex.Location{0, 4}},
		},
		Error: io.EOF,
	},
	lexData{
		// "\`" -- checking an escaped backtick within a double-quoted string
		Input: "\"\\`\"",
		Tokens: []*lex.Token{
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 0}, lex.Location{0, 1}},
			&lex.Token{lex.StringSegment, "`", lex.Location{0, 1}, lex.Location{0, 3}},
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 3}, lex.Location{0, 4}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: `"\""`,
		Tokens: []*lex.Token{
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 0}, lex.Location{0, 1}},
			&lex.Token{lex.StringSegment, `"`, lex.Location{0, 1}, lex.Location{0, 3}},
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 3}, lex.Location{0, 4}},
		},
		Error: io.EOF,
	},
	lexData{
		Input: `"$\$$\$"`,
		Tokens: []*lex.Token{
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 0}, lex.Location{0, 1}},
			// Is this what we want? Escaped dollars cannot be used to start
			// variables or param expansions, so this is probably fine.
			&lex.Token{lex.Dollar, `$`, lex.Location{0, 1}, lex.Location{0, 2}},
			&lex.Token{lex.StringSegment, `$`, lex.Location{0, 2}, lex.Location{0, 4}},
			&lex.Token{lex.Dollar, `$`, lex.Location{0, 4}, lex.Location{0, 5}},
			&lex.Token{lex.StringSegment, `$`, lex.Location{0, 5}, lex.Location{0, 7}},
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 7}, lex.Location{0, 8}},
		},
		Error: io.EOF,
	},
	lexData{
		Input:  `"unclosed-string`,
		Tokens: nil,
		Error:  fmt.Errorf("Unclosed double-quoted string"),
	},
	lexData{
		Input: `"${WUMBO}"`,
		Tokens: []*lex.Token{
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 0}, lex.Location{0, 1}},
			&lex.Token{lex.Dollar, `$`, lex.Location{0, 1}, lex.Location{0, 2}},
			&lex.Token{lex.LeftBrace, `{`, lex.Location{0, 2}, lex.Location{0, 3}},
			&lex.Token{lex.Name, `WUMBO`, lex.Location{0, 3}, lex.Location{0, 8}},
			&lex.Token{lex.RightBrace, `}`, lex.Location{0, 8}, lex.Location{0, 9}},
			&lex.Token{lex.DoubleQuote, `"`, lex.Location{0, 9}, lex.Location{0, 10}},
		},
		Error: io.EOF,
	},
}
