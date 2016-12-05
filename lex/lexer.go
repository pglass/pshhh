package lex

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

const eof = -1

type Lexer struct {
	input  string // the string we are scanning
	start  int    // the start of the next token to emit
	pos    int    // the current position
	line   int
	tokens chan Token
	state  stateFn

	peekBuf []Token
}

func NewLexer(input string) *Lexer {
	lexer := &Lexer{
		input:   input,
		line:    1,
		tokens:  make(chan Token),
		peekBuf: []Token{},
	}
	go lexer.run()
	return lexer
}

func (lx *Lexer) nextRune() rune {
	c, size := lx.peek()
	lx.pos += size
	if c == '\n' {
		lx.line += 1
	}
	return c
}

func (lx *Lexer) peekRune() rune {
	c, _ := lx.peek()
	return c
}

func (lx *Lexer) peek() (rune, int) {
	if int(lx.pos) >= len(lx.input) {
		return eof, 0
	}
	return utf8.DecodeRuneInString(lx.input[lx.pos:])
}

func (lx *Lexer) Next() Token {
	if len(lx.peekBuf) > 0 {
		item := lx.peekBuf[0]
		lx.peekBuf = lx.peekBuf[1:]
		return item
	} else {
		return <-lx.tokens
	}
}

func (lx *Lexer) Peek() Token {
	if len(lx.peekBuf) == 0 {
		lx.peekBuf = append(lx.peekBuf, <-lx.tokens)
	}
	return lx.peekBuf[0]
}

func (lx *Lexer) Unread(tok Token) {
	lx.peekBuf = append([]Token{tok}, lx.peekBuf...)
}

func (lx *Lexer) HasAnyToken(ttypes ...TokenType) bool {
	tok := lx.Peek()
	for _, ttype := range ttypes {
		if tok.Type == ttype {
			return true
		}
	}
	return false
}

// emit the current token to the token channel
func (lx *Lexer) emit(ttype TokenType) {
	lx.emitText(ttype, lx.input[lx.start:lx.pos])
}

func (lx *Lexer) emitText(ttype TokenType, text string) {
	token := Token{
		Type: ttype,
		Text: text,
		Pos:  lx.start,
		Line: lx.line,
	}

	token.Line -= strings.Count(token.Text, "\n")

	switch token.Type {
	case Name:
		if ttype, ok := RESERVED_WORDS[token.Text]; ok {
			token.Type = ttype
		}
	}
	lx.tokens <- token
	lx.start = lx.pos

}

func (lx *Lexer) emitBuffer(ttype TokenType, buffer bytes.Buffer) {
	if buffer.Len() > 0 {
		lx.emitText(ttype, buffer.String())
	}
	// always shift the start pos when we emit
	lx.start = lx.pos
}

// emit an error token to the token channel
func (lx *Lexer) errorf(format string, args ...interface{}) stateFn {
	lx.tokens <- Token{
		Type: ERROR,
		Text: fmt.Sprintf(format, args...),
		Pos:  lx.start,
		Line: lx.line,
	}
	return nil
}

// this is run in a coroutine. it calls each state function in succession.
// the state functions use the lexer to emit tokens, by calling lx.emit().
func (lx *Lexer) run() {
	for lx.state = lexText; lx.state != nil; {
		lx.state = lx.state(lx, lexText)
	}
	close(lx.tokens)
}

func (lx *Lexer) hasString(s string) bool {
	return strings.HasPrefix(lx.input[lx.pos:], s)
}
