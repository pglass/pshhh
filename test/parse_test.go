package test

import (
	"testing"

	"github.com/pglass/pshhh/ast"
	"github.com/pglass/pshhh/lex"
	"github.com/stretchr/testify/assert"
)

func parse_input(input string) (ast.Node, error) {
	lexer := lex.NewLexer(input)
	parser := ast.NewParser(lexer)
	return parser.Parse()
}

type parseData struct {
	Input  string
	Output ast.Node
	Error  error
}

func run_parser_test(t *testing.T, data parseData) {
	t.Logf("Input: %q", data.Input)
	t.Logf("Error (expected): %v", data.Error)
	t.Logf("Output (expected): %v", data.Output)

	node, err := parse_input(data.Input)

	t.Logf("Error (received): %v", err)
	t.Logf("Output (received): %v", node)

	if data.Error != nil {
		assert.Equal(t, data.Error, err)
	} else {
		assert.Nil(t, err)
	}
	assert.Equal(t, data.Output, node)
}

func TestParser(t *testing.T) {
	for _, data := range PARSE_CASES {
		t.Run(data.Input, func(t *testing.T) { run_parser_test(t, data) })
	}
}

var PARSE_CASES = []parseData{
	parseData{
		Input: "echo",
		Output: ast.NewGenericNode(
			&ast.CommandList{
				Separators: []lex.Token{},
				Commands: []ast.Command{
					&ast.SimpleCommand{
						Redirects: []*ast.IoRedirect{},
						Words: []*ast.Str{
							ast.NewStrFromTok(lex.Token{lex.Name, "echo", 0, 1}),
						},
					},
				},
			},
		),
		Error: nil,
	},
	parseData{
		Input: "echo 1 2; echo wumbo a & echo;",
		Output: ast.NewGenericNode(
			&ast.CommandList{
				Separators: []lex.Token{
					lex.Token{lex.Semi, ";", 8, 1},
					lex.Token{lex.Ampersand, "&", 23, 1},
					lex.Token{lex.Semi, ";", 29, 1},
				},
				Commands: []ast.Command{
					&ast.SimpleCommand{
						Redirects: []*ast.IoRedirect{},
						Words: []*ast.Str{
							ast.NewStrFromTok(lex.Token{lex.Name, "echo", 0, 1}),
							ast.NewStrFromTok(lex.Token{lex.Number, "1", 5, 1}),
							ast.NewStrFromTok(lex.Token{lex.Number, "2", 7, 1}),
						},
					},
					&ast.SimpleCommand{
						Redirects: []*ast.IoRedirect{},
						Words: []*ast.Str{
							ast.NewStrFromTok(lex.Token{lex.Name, "echo", 10, 1}),
							ast.NewStrFromTok(lex.Token{lex.Name, "wumbo", 15, 1}),
							ast.NewStrFromTok(lex.Token{lex.Name, "a", 21, 1}),
						},
					},
					&ast.SimpleCommand{
						Redirects: []*ast.IoRedirect{},
						Words: []*ast.Str{
							ast.NewStrFromTok(lex.Token{lex.Name, "echo", 25, 1}),
						},
					},
				},
			},
		),
		Error: nil,
	},
	parseData{
		Input: `"$MINI"`,
		Output: ast.NewGenericNode(
			&ast.Str{
				Pieces: []ast.StrPiece{
					&ast.ParameterExpansion{
						VarName:  &lex.Token{lex.Name, "MINI", 2, 1},
						Operator: nil,
						Word:     nil,
					},
				},
			},
		),
		Error: nil,
	},
	parseData{
		Input: `"${MINI}"`,
		Output: ast.NewGenericNode(
			&ast.Str{
				Pieces: []ast.StrPiece{
					&ast.ParameterExpansion{
						VarName:  &lex.Token{lex.Name, "MINI", 3, 1},
						Operator: nil,
						Word:     nil,
					},
				},
			},
		),
		Error: nil,
	},
}
