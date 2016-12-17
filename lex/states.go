package lex

import (
	"bytes"
	"strings"
	"unicode"
)

type stateFn func(*Lexer, stateFn) stateFn

func composeStates(lx *Lexer, state stateFn, states ...stateFn) stateFn {
	nextState := lexText
	if len(states) > 1 {
		nextState = composeStates(lx, states[0], states[1:]...)
	} else if len(states) == 1 {
		nextState = states[0]
	}
	return func(*Lexer, stateFn) stateFn {
		return state(lx, nextState)
	}
}

func lexText(lx *Lexer, nextState stateFn) stateFn {
	c := lx.peekRune()
	if c == eof {
		lx.emit(EOF)
		return nil
	} else if c == '$' {
		return lexDollarExpansion(lx, nextState)
	} else if IsWordChar(c) {
		return lexWord(lx, nextState)
	} else if unicode.IsSpace(c) {
		return lexSpace(lx, nextState)
	} else if c == '\'' {
		return lexSingleQuotedString(lx, nextState)
	} else if c == '"' {
		return lexDoubleQuotedString(lx, nextState)
	} else if unicode.IsPunct(c) || unicode.IsSymbol(c) {
		return lexOperator(lx, nextState)
	} else {
		return lx.errorf("Unexpected rune %q", c)
	}
}

// Read a Word token (which may actually be a Name or Number)
func lexWord(lx *Lexer, nextState stateFn) stateFn {
	// A WORD consisting of only digits is a NUMBER
	// A WORD that does not start with a digit is a NAME
	if unicode.IsDigit(lx.peekRune()) {
		return lexNumberOrWord(lx, nextState)
	} else {
		for IsWordChar(lx.peekRune()) {
			lx.nextRune()
		}
	}

	lx.emit(Name)
	return nextState
}

func lexNumberOrWord(lx *Lexer, nextState stateFn) stateFn {
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
	return nextState
}

func lexSpace(lx *Lexer, nextState stateFn) stateFn {
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
	return nextState
}

func lexOperator(lx *Lexer, nextState stateFn) stateFn {
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
	return nextState
}

func lexSingleQuotedString(lx *Lexer, nextState stateFn) stateFn {
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
	return nextState
}

func lexDoubleQuotedString(lx *Lexer, nextState stateFn) stateFn {
	if c := lx.nextRune(); c != '"' {
		return lx.errorf("Expected double quote to start a string (got %c)", c)
	} else {
		lx.emit(DoubleQuote)
	}

	return composeStates(lx, lexDoubleQuotedStringContents, nextState)
}

func lexDoubleQuotedStringContents(lx *Lexer, nextState stateFn) stateFn {
	var buffer bytes.Buffer
	for {
		c := lx.peekRune()
		if c == '"' {
			lx.emitBuffer(StringSegment, buffer)
			lx.nextRune()
			lx.emit(DoubleQuote)
			break
		} else if c == eof {
			return lx.errorf("Unclosed string")
		} else if c == '\\' {
			lx.nextRune()
			cc := lx.peekRune()
			if strings.ContainsRune("$`\"\\", cc) {
				lx.nextRune()
				buffer.WriteRune(cc)
			} else if cc == '\n' {
				lx.nextRune()
			} else {
				buffer.WriteRune(c)
			}
		} else if c == '$' {
			lx.emitBuffer(StringSegment, buffer)
			return composeStates(lx, lexDollarExpansion, lexDoubleQuotedStringContents, nextState)
		} else {
			lx.nextRune()
			buffer.WriteRune(c)
		}
	}
	return nextState
}

func lexDollarExpansion(lx *Lexer, nextState stateFn) stateFn {
	if c := lx.nextRune(); c != '$' {
		return lx.errorf("Expected '$' to start dollar expansion (got %c)", c)
	} else {
		lx.emit(Dollar)
	}

	c := lx.peekRune()
	if c == '{' {
		return lexBraceExpansion(lx, nextState)
	} else if c == '(' {
		return lexParenExpansion(lx, nextState)
	} else if IsNameChar(c) {
		return lexName(lx, nextState)
	}
	return nextState
}

func lexBraceExpansion(lx *Lexer, nextState stateFn) stateFn {
	if c := lx.nextRune(); c != '{' {
		return lx.errorf("Expected '{' to start brace expansion (got %c)", c)
	} else {
		lx.emit(LeftBrace)
	}

	c := lx.peekRune()
	if IsNameChar(c) {
		return composeStates(lx, lexName, lexBraceExpansionEnd, nextState)
	}

	return nextState
}

func lexName(lx *Lexer, nextState stateFn) stateFn {
	for c := lx.peekRune(); IsNameChar(c); c = lx.peekRune() {
		lx.nextRune()
	}
	if lx.pos <= lx.start {
		return lx.errorf("Expected Name in brace expansion")
	}
	lx.emit(Name)
	return nextState
}

func lexBraceExpansionEnd(lx *Lexer, nextState stateFn) stateFn {
	c := lx.peekRune()
	if strings.ContainsRune(":+-=?", c) {
		return composeStates(lx, lexOperator, lexText, lexOperator, nextState)
	} else if c == '}' {
		lx.nextRune()
		lx.emit(RightBrace)
	} else {
		return lx.errorf("Unclosed brace expansion (expected '}')")
	}
	return nextState
}

func lexParenExpansion(lx *Lexer, nextState stateFn) stateFn {
	return nil
}
