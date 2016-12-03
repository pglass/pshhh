package test

import (
	"github.com/stretchr/testify/assert"
	"psh/ast"
	"psh/lex"
	// "reflect"
	"strings"
	"testing"
)

func parse_input(input string) (ast.Node, error) {
	reader := strings.NewReader(input)
	lexer := lex.NewLexer(reader)
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
				Separators: []*lex.Token{},
				Commands: []ast.Command{
					&ast.SimpleCommand{
						Redirects: []*ast.IoRedirect{},
						Words: []*ast.String{
							ast.NewStringFromTok(&lex.Token{lex.Name, "echo", lex.Location{0, 0}, lex.Location{0, 4}}),
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
				Separators: []*lex.Token{
					&lex.Token{lex.Semi, ";", lex.Location{0, 8}, lex.Location{0, 9}},
					&lex.Token{lex.Ampersand, "&", lex.Location{0, 23}, lex.Location{0, 24}},
					&lex.Token{lex.Semi, ";", lex.Location{0, 29}, lex.Location{0, 30}},
				},
				Commands: []ast.Command{
					&ast.SimpleCommand{
						Redirects: []*ast.IoRedirect{},
						Words: []*ast.String{
							ast.NewStringFromTok(&lex.Token{lex.Name, "echo", lex.Location{0, 0}, lex.Location{0, 4}}),
							ast.NewStringFromTok(&lex.Token{lex.Number, "1", lex.Location{0, 5}, lex.Location{0, 6}}),
							ast.NewStringFromTok(&lex.Token{lex.Number, "2", lex.Location{0, 7}, lex.Location{0, 8}}),
						},
					},
					&ast.SimpleCommand{
						Redirects: []*ast.IoRedirect{},
						Words: []*ast.String{
							ast.NewStringFromTok(&lex.Token{lex.Name, "echo", lex.Location{0, 10}, lex.Location{0, 14}}),
							ast.NewStringFromTok(&lex.Token{lex.Name, "wumbo", lex.Location{0, 15}, lex.Location{0, 20}}),
							ast.NewStringFromTok(&lex.Token{lex.Name, "a", lex.Location{0, 21}, lex.Location{0, 22}}),
						},
					},
					&ast.SimpleCommand{
						Redirects: []*ast.IoRedirect{},
						Words: []*ast.String{
							ast.NewStringFromTok(&lex.Token{lex.Name, "echo", lex.Location{0, 25}, lex.Location{0, 29}}),
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
			&ast.String{
				Pieces: []ast.StringPiece{
					&ast.ParameterExpansion{
						VarName:  &lex.Token{lex.Name, "MINI", lex.Location{0, 2}, lex.Location{0, 6}},
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
			&ast.String{
				Pieces: []ast.StringPiece{
					&ast.ParameterExpansion{
						VarName:  &lex.Token{lex.Name, "MINI", lex.Location{0, 3}, lex.Location{0, 7}},
						Operator: nil,
						Word:     nil,
					},
				},
			},
		),
		Error: nil,
	},
}
