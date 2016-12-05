package lex

import (
	"fmt"
	"strings"
	"unicode"
)

func IsWordChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || strings.ContainsRune("./=-", c)
}

func IsNameChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}

type Token struct {
	Type TokenType
	Text string
	Pos  int
	Line int
}

func (t Token) Format(f fmt.State, c rune) {
	if t.Type == ERROR {
		fmt.Fprintf(f, "%v (line %v): %v", t.Type, t.Line, t.Text)
	} else {
		fmt.Fprintf(f, "%v(%q %v, %v)", t.Type, t.Text, t.Pos, t.Line)
	}
}

// When adding a token type here, update the TokenTypeString and
// TokenTypeFromString functions as appropriate.
type TokenType int

const (
	Unknown TokenType = iota
	EOF
	ERROR
	Word
	AssignmentWord
	Name
	Newline
	Number
	Space

	Ampersand
	AndIf
	Pipe
	OrIf
	Semi
	DoubleSemi

	Less
	DoubleLess
	Great
	DoubleGreat
	LessAnd
	GreatAnd
	LessGreat
	DoubleLessDash
	Clobber
	Dollar
	SingleQuote
	DoubleQuote
	StringSegment

	If
	Then
	Else
	Elif
	Fi
	Do
	Done

	Case
	Esac
	While
	Until
	For
	Function

	LeftBrace
	RightBrace
	LeftParen
	RightParen
	Bang
	In

	ColonDash
	ColonQuestion
	ColonPlus
	ColonEquals
	Plus
	Dash
	Equals
	Question
)

var OPERATORS = map[string]TokenType{
	"&":   Ampersand,
	"&&":  AndIf,
	"|":   Pipe,
	"||":  OrIf,
	";":   Semi,
	";;":  DoubleSemi,
	"<":   Less,
	">":   Great,
	"<<":  DoubleLess,
	">>":  DoubleGreat,
	"<&":  LessAnd,
	">&":  GreatAnd,
	"<>":  LessGreat,
	"<<-": DoubleLessDash,
	">|":  Clobber,
	"$":   Dollar,
	"{":   LeftBrace,
	"}":   RightBrace,
	"(":   LeftParen,
	")":   RightParen,
	"!":   Bang,

	":-": ColonDash,
	":=": ColonEquals,
	":?": ColonQuestion,
	":+": ColonPlus,
	"-":  Dash,
	"=":  Equals,
	"?":  Question,
	"+":  Plus,
}

// e.g. the ':=' in ${P:=W}
var PARAMETER_EXPANSION_TTYPES = []TokenType{
	ColonDash, ColonEquals, ColonQuestion, ColonPlus,
	Dash, Equals, Question, Plus,
}

var RESERVED_WORDS = map[string]TokenType{
	"if":       If,
	"then":     Then,
	"else":     Else,
	"elif":     Elif,
	"fi":       Fi,
	"do":       Do,
	"done":     Done,
	"case":     Case,
	"esac":     Esac,
	"while":    While,
	"until":    Until,
	"for":      For,
	"function": Function,
	"in":       In,
}

var tokenToString = map[TokenType]string{
	Unknown:        "Unkown",
	EOF:            "EOF",
	ERROR:          "ERROR",
	Word:           "Word",
	AssignmentWord: "AssignmentWord",
	Name:           "Name",
	Newline:        "Newline",
	Number:         "Number",
	Space:          "Space",
	Ampersand:      "Ampersand",
	AndIf:          "AndIf",
	Pipe:           "Pipe",
	OrIf:           "OrIf",
	Semi:           "Semi",
	DoubleSemi:     "DoubleSemi",
	Less:           "Less",
	DoubleLess:     "DoubleLess",
	Great:          "Great",
	DoubleGreat:    "DoubleGreat",
	LessAnd:        "LessAnd",
	GreatAnd:       "GreatAnd",
	LessGreat:      "LessGreat",
	DoubleLessDash: "DoubleLessDash",
	Clobber:        "Clobber",
	Dollar:         "Dollar",
	SingleQuote:    "SingleQuote",
	DoubleQuote:    "DoubleQuote",
	StringSegment:  "StringSegment",

	If:       "If",
	Then:     "Then",
	Else:     "Else",
	Elif:     "Elif",
	Fi:       "Fi",
	Do:       "Do",
	Done:     "Done",
	Case:     "Case",
	Esac:     "Esac",
	While:    "While",
	Until:    "Until",
	For:      "For",
	Function: "Function",

	LeftBrace:  "LeftBrace",
	RightBrace: "RightBrace",
	LeftParen:  "LeftParen",
	RightParen: "RightParen",
	Bang:       "Bang",
	In:         "In",

	ColonDash:     "ColonDash",
	ColonQuestion: "ColonQuestion",
	ColonPlus:     "ColonPlus",
	ColonEquals:   "ColonEquals",
	Plus:          "Plus",
	Dash:          "Dash",
	Equals:        "Equals",
	Question:      "Question",
}

func (tt TokenType) Format(f fmt.State, c rune) {
	if text, ok := tokenToString[tt]; ok {
		fmt.Fprintf(f, "%v", text)
	} else {
		fmt.Fprintf(f, "<TokenType %v>", tt)
	}
}
