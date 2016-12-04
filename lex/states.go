package lex

import (
	"unicode"
)

type stateFn func(*Lexer) stateFn

func lexText(lx *Lexer) stateFn {
	c := lx.peekRune()
	if c == eof {
		lx.emit(EOF)
		return nil
	} else if IsWordChar(c) {
		return lexWord
	} else if unicode.IsSpace(c) {
		return lexSpace
	} else if c == '\'' {
		return lexSingleQuotedString
	} else if unicode.IsPunct(c) || unicode.IsSymbol(c) {
		return lexOperator
	} else {
		return lx.errorf("Unexpected rune %q", c)
	}
}

// Read a Word token (which may actually be a Name or Number)
func lexWord(lx *Lexer) stateFn {
	// A WORD consisting of only digits is a NUMBER
	// A WORD that does not start with a digit is a NAME
	if unicode.IsDigit(lx.peekRune()) {
		return lexNumberOrWord
	}

	for IsWordChar(lx.peekRune()) {
		lx.nextRune()
	}

	lx.emit(Name)
	return lexText
}

// Read a Number or Word token
func lexNumberOrWord(lx *Lexer) stateFn {
	if c := lx.nextRune(); !unicode.IsDigit(c) {
		return lx.errorf("Expected Name or Number to start with a digit (got %c)", c)
	}

	ttype := Number
	for IsWordChar(lx.peekRune()) {
		c := lx.nextRune()
		if !unicode.IsDigit(c) {
			ttype = Word
		}
	}
	lx.emit(ttype)
	return lexText
}

func lexSpace(lx *Lexer) stateFn {
	if c := lx.nextRune(); !unicode.IsSpace(c) {
		return lx.errorf("Expected Space or Newline to start with a space char (got %c)", c)
	} else if c == '\n' {
		lx.emit(Newline)
	} else {
		for {
			c := lx.peekRune()
			// emit newlines in separate tokens
			if c == '\n' {
				break
			} else if unicode.IsSpace(c) {
				lx.nextRune()
			} else {
				break
			}
		}
		lx.emit(Space)
	}
	return lexText
}

func lexOperator(lx *Lexer) stateFn {
	var result_op string = ""
	var result_ttype TokenType

	// consume the longest matching operator
	for op, ttype := range OPERATORS {
		if lx.hasString(op) && len(op) > len(result_op) {
			result_op = op
			result_ttype = ttype
		}
	}

	if len(result_op) == 0 {
		return lx.errorf("Failed to read operator (saw char %c)", lx.peekRune())
	}

	// consume the text
	for range result_op {
		lx.nextRune()
	}
	lx.emit(result_ttype)
	return lexText
}

func lexSingleQuotedString(lx *Lexer) stateFn {
	if c := lx.nextRune(); c != '\'' {
		return lx.errorf("Expected single quote to start a string (got %c)", c)
	} else {
		lx.emit(SingleQuote)
	}

	for lx.peekRune() != '\'' {
		c := lx.nextRune()
		if c == eof {
			break
		}
	}

	if c := lx.peekRune(); c != '\'' {
		return lx.errorf("Unclosed string '%v...", lx.input[lx.start:lx.start+10])
	}

	lx.emit(StringSegment)
	lx.nextRune()
	lx.emit(SingleQuote)
	return lexText
}
