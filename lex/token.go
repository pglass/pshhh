package lex

import (
	"fmt"
	"strings"
	"unicode"
)

func IsWordChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || strings.ContainsRune("./=-", c)
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
	Equals

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
)

func (tt TokenType) Format(f fmt.State, c rune) {
	text := ""
	switch tt {
	case Unknown:
		text = "Unkown"
	case EOF:
		text = "EOF"
	case ERROR:
		text = "ERROR"
	case Word:
		text = "Word"
	case AssignmentWord:
		text = "AssignmentWord"
	case Name:
		text = "Name"
	case Newline:
		text = "Newline"
	case Number:
		text = "Number"
	case Space:
		text = "Space"

	case Ampersand:
		text = "Ampersand"
	case AndIf:
		text = "AndIf"
	case Pipe:
		text = "Pipe"
	case OrIf:
		text = "OrIf"
	case Semi:
		text = "Semi"
	case DoubleSemi:
		text = "DoubleSemi"

	case Less:
		text = "Less"
	case DoubleLess:
		text = "DoubleLess"
	case Great:
		text = "Great"
	case DoubleGreat:
		text = "DoubleGreat"
	case LessAnd:
		text = "LessAnd"
	case GreatAnd:
		text = "GreatAnd"
	case LessGreat:
		text = "LessGreat"
	case DoubleLessDash:
		text = "DoubleLessDash"
	case Clobber:
		text = "Clobber"
	case Dollar:
		text = "Dollar"
	case SingleQuote:
		text = "SingleQuote"
	case DoubleQuote:
		text = "DoubleQuote"
	case StringSegment:
		text = "StringSegment"

	case If:
		text = "If"
	case Then:
		text = "Then"
	case Else:
		text = "Else"
	case Elif:
		text = "Elif"
	case Fi:
		text = "Fi"
	case Do:
		text = "Do"
	case Done:
		text = "Done"

	case Case:
		text = "Case"
	case Esac:
		text = "Esac"
	case While:
		text = "While"
	case Until:
		text = "Until"
	case For:
		text = "For"
	case Function:
		text = "Function"

	case LeftBrace:
		text = "LeftBrace"
	case RightBrace:
		text = "RightBrace"
	case LeftParen:
		text = "LeftParen"
	case RightParen:
		text = "RightParen"
	case Bang:
		text = "Bang"
	case In:
		text = "In"
	case Equals:
		text = "Equals"

	default:
		text = fmt.Sprintf("<TokenType %v>", tt)
	}
	fmt.Fprintf(f, "%v", text)
}

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
	"'":   SingleQuote,
	`"`:   DoubleQuote,
	"{":   LeftBrace,
	"}":   RightBrace,
	"(":   LeftParen,
	")":   RightParen,
	"!":   Bang,
	"=":   Equals,
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

// func TokenTypeFromString(text string) TokenType {
// 	switch text {
// 	case "&":
// 		return Ampersand
// 	case "&&":
// 		return AndIf
// 	case "|":
// 		return Pipe
// 	case "||":
// 		return OrIf
// 	case ";":
// 		return Semi
// 	case ";;":
// 		return DoubleSemi
//
// 	case "<":
// 		return Less
// 	case ">":
// 		return Great
// 	case "<<":
// 		return DoubleLess
// 	case ">>":
// 		return DoubleGreat
// 	case "<&":
// 		return LessAnd
// 	case ">&":
// 		return GreatAnd
// 	case "<>":
// 		return LessGreat
// 	case "<<-":
// 		return DoubleLessDash
// 	case ">|":
// 		return Clobber
// 	case "$":
// 		return Dollar
// 	case "'":
// 		return SingleQuote
// 	case `"`:
// 		return DoubleQuote
//
// 	case "if":
// 		return If
// 	case "then":
// 		return Then
// 	case "else":
// 		return Else
// 	case "elif":
// 		return Elif
// 	case "fi":
// 		return Fi
// 	case "do":
// 		return Do
// 	case "done":
// 		return Done
//
// 	case "case":
// 		return Case
// 	case "esac":
// 		return Esac
// 	case "while":
// 		return While
// 	case "until":
// 		return Until
// 	case "for":
// 		return For
// 	case "function":
// 		return Function
//
// 	case "{":
// 		return LeftBrace
// 	case "}":
// 		return RightBrace
// 	case "(":
// 		return LeftParen
// 	case ")":
// 		return RightParen
// 	case "!":
// 		return Bang
// 	case "in":
// 		return In
// 	}
//
// 	return Unknown
// }
//

func possiblePunctation(c rune) bool {
	switch c {
	case ';':
		return true
	case '$':
		return true
	case '&':
		return true
	case '|':
		return true
	case '<':
		return true
	case '>':
		return true
	case '/':
		return true
	case '{':
		return true
	case '}':
		return true
	case '!':
		return true
	case '(':
		return true
	case ')':
		return true
	case '=':
		return true
	}
	return false
}
