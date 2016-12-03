package lex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"unicode"
)

type Lexer struct {
	reader   *bufio.Reader
	peekbuf  []*peekbufEntry
	location Location
}

type peekbufEntry struct {
	Token *Token
	Error error
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader:   bufio.NewReaderSize(reader, 2),
		peekbuf:  []*peekbufEntry{},
		location: Location{0, 0},
	}
}

func (lx *Lexer) HasAnyToken(ttypes ...TokenType) bool {
	if tok, err := lx.Peek(); err == nil {
		for _, ttype := range ttypes {
			if tok.Type == ttype {
				return true
			}
		}
	}
	return false
}

func (lx *Lexer) Peek() (*Token, error) {
	if len(lx.peekbuf) == 0 {
		err := lx.nextToken()
		if err != nil {
			lx.peekbuf = append(lx.peekbuf, &peekbufEntry{nil, err})
		}
	}
	if len(lx.peekbuf) > 0 {
		entry := lx.peekbuf[0]
		return entry.Token, entry.Error
	}
	return nil, io.EOF
}

func (lx *Lexer) Unread(tok *Token) {
	lx.peekbuf = append(lx.peekbuf, &peekbufEntry{tok, nil})
}

/* Read the next token. Returns a non-nil error on syntax errors. The error
 * is io.EOF on EOF. */
func (lx *Lexer) Next() (*Token, error) {
	if len(lx.peekbuf) == 0 {
		err := lx.nextToken()
		if err != nil {
			lx.peekbuf = append(lx.peekbuf, &peekbufEntry{nil, err})
		}
	}
	if len(lx.peekbuf) > 0 {
		entry := lx.peekbuf[0]
		lx.peekbuf = lx.peekbuf[1:]
		log.Printf("Token = %v", entry.Token)
		return entry.Token, entry.Error
	}
	return nil, io.EOF
}

func (lx *Lexer) nextToken() error {
	c, _, err := lx.peekRune()
	if err != nil {
		return err
	}

	if IsWordChar(c) {
		err = lx.readAlphaNumeric(c)
	} else if unicode.IsSpace(c) {
		err = lx.readSpace(c)
	} else if c == '\'' {
		err = lx.readSingleQuotedString(c)
	} else if c == '"' {
		err = lx.readDoubleQuotedString(c)
	} else if c == '$' {
		err = lx.readDollarExpansion(c)
	} else if lx.possiblePunctation(c) {
		err = lx.readOperator(c)
	} else {
		c, _, err = lx.readRune()
		log.Printf("Unhandled rune '%c'", c)
		if err != io.EOF {
			err = fmt.Errorf("Syntax error near %q", c)
		}
	}

	return err
}

func (lx *Lexer) peekRune() (rune, int, error) {
	c, size, err := lx.reader.ReadRune()
	if err == nil {
		err = lx.reader.UnreadRune()
	}
	return c, size, err
}

func (lx *Lexer) hasRune(wants ...rune) bool {
	c, _, err := lx.peekRune()
	if err == nil {
		for _, want := range wants {
			if c == want {
				return true
			}
		}
	}
	return false
}

func (lx *Lexer) readRune() (rune, int, error) {
	c, size, err := lx.reader.ReadRune()
	if err == nil {
		if c == '\n' {
			lx.location.Column = 0
			lx.location.LineNumber += 1
		} else {
			lx.location.Column += 1
		}
	}
	return c, size, err
}

func (lx *Lexer) yieldToken(token *Token) {
	lx.peekbuf = append(lx.peekbuf, &peekbufEntry{token, nil})
}

func (lx *Lexer) readAlphaNumeric(peekChar rune) error {
	start := lx.location.Copy()

	var buffer bytes.Buffer
	for {
		c, _, err := lx.peekRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if IsWordChar(c) {
			lx.readRune()
		} else {
			break
		}
		buffer.WriteRune(c)
	}

	token := &Token{Word, buffer.String(), start, lx.location.Copy()}

	// A WORD that does not start with a digit is a NAME
	// A WORD consisting of only digits is a NUMBER
	ttype := Number
	for _, c := range token.Text {
		if unicode.IsDigit(c) {
			continue
		}

		if !unicode.IsDigit(peekChar) {
			ttype = Name
		} else {
			ttype = Word
		}
		break
	}
	token.Type = ttype

	// convert "if", "else" to If, Else tokens
	if ttype, match := lx.detectReservedWord(token); match {
		token.Type = ttype
	}

	lx.yieldToken(token)

	return nil
}

func (lx *Lexer) readSpace(peekChar rune) error {
	if peekChar == '\n' {
		lx.yieldNextRuneAsToken(Newline)
		return nil
	}

	// output all adjacent spaces (not newlines) in a single token
	start := lx.location.Copy()
	var buffer bytes.Buffer
	for {
		c, _, err := lx.peekRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if c == '\n' {
			break
		} else if unicode.IsSpace(c) {
			lx.readRune()
			buffer.WriteRune(c)
		} else {
			break
		}
	}
	token := &Token{Space, buffer.String(), start, lx.location.Copy()}
	lx.yieldToken(token)
	return nil
}

func (lx *Lexer) readOperator(peekChar rune) error {
	start := lx.location.Copy()

	c, _, err := lx.readRune()
	if err != nil {
		return err
	}

	// Return text or (text + c) depending if the next char is c
	maybeAppendNextChar := func(text string, c rune) string {
		if lx.hasRune(c) {
			lx.readRune()
			return fmt.Sprintf("%s%c", text, c)
		} else {
			return text
		}
	}

	text := ""
	switch c {
	case '&':
		text = maybeAppendNextChar("&", '&')
	case '|':
		text = maybeAppendNextChar("|", '|')
	case ';':
		text = maybeAppendNextChar(";", ';')
	case '<':
		text = maybeAppendNextChar("<", '<')
		if text == "<<" && lx.hasRune('-') {
			text = maybeAppendNextChar("<<", '-')
		} else if text == "<" && lx.hasRune('&') {
			text = maybeAppendNextChar("<", '&')
		} else if text == "<" && lx.hasRune('>') {
			text = maybeAppendNextChar("<", '>')
		}
	case '>':
		if lx.hasRune('>') {
			text = maybeAppendNextChar(">", '>')
		} else if lx.hasRune('&') {
			text = maybeAppendNextChar(">", '&')
		} else if lx.hasRune('|') {
			text = maybeAppendNextChar(">", '}')
		} else {
			text = ">"
		}
	case '$':
		text = "$"
	default:
		return fmt.Errorf("Invalid operator found near %c", c)
	}

	ttype := TokenTypeFromString(text)
	if ttype == Unknown {
		if text != "" {
			return fmt.Errorf("Invalid token near %q", text)
		} else {
			return fmt.Errorf("Syntax error")
		}
	}

	token := &Token{ttype, text, start, lx.location.Copy()}
	lx.yieldToken(token)

	return nil
}

func (lx *Lexer) readSingleQuotedString(peekChar rune) error {
	start := lx.location.Copy()

	if peekChar != '\'' {
		return fmt.Errorf("Invalid usage of readSingleQuotedString [bug!]")
	} else {
		lx.readRune()
	}

	var buffer bytes.Buffer
	for {
		c, _, err := lx.peekRune()
		if err == io.EOF {
			return fmt.Errorf("Unclosed string '%v...", buffer.String()[:10])
		} else if err != nil {
			return err
		} else if c == '\'' {
			lx.readRune()
			break
		} else {
			buffer.WriteRune(c)
			lx.readRune()
		}
	}

	token := &Token{StringSegment, buffer.String(), start, lx.location.Copy()}
	lx.yieldToken(token)

	return nil
}

/* This will yield tokens directly to peekbuf. If an error is returned, the
 * lexer will be in an inconsistent state (consumed chars from the reader,
 * peekbuf full of tokens, etc).
 */
func (lx *Lexer) readDoubleQuotedString(peekChar rune) error {
	if peekChar != '"' {
		return fmt.Errorf("Invalid usage of readDoubleQuotedString [bug!]")
	} else {
		lx.yieldNextRuneAsToken(DoubleQuote)
	}

	start := lx.location.Copy()
	var buffer bytes.Buffer
	for {
		c, _, err := lx.peekRune()
		if err == io.EOF {
			return fmt.Errorf(`Unclosed double-quoted string`)
		} else if err != nil {
			return err
		}

		switch c {
		case '"':
			// yield the current buffer as a token
			text := buffer.String()
			if len(text) > 0 {
				token := &Token{StringSegment, text, start, lx.location.Copy()}
				lx.yieldToken(token)
				start = lx.location.Copy()
				buffer.Reset()
			}

			// yield the double quote as a token
			lx.yieldNextRuneAsToken(DoubleQuote)
			return nil
		case '\\':
			lx.readRune()
			// support blackslash escapes
			if lx.hasRune('$', '`', '"', '\\') {
				cc, _, _ := lx.readRune()
				buffer.WriteRune(cc)
			} else {
				buffer.WriteRune(c)
			}
		case '$':
			// input like "abc${WUMBO}def" produces multiple tokens:
			//	  StringSegment "abc"
			//	  <parameter expansion tokens>
			//	  StringSegment "def"
			text := buffer.String()
			if len(text) > 0 {
				token := &Token{StringSegment, text, start, lx.location.Copy()}
				lx.yieldToken(token)
				buffer.Reset()
			}

			err := lx.readDollarExpansion(c)
			if err != nil {
				return err
			}
			start = lx.location.Copy()
		default:
			buffer.WriteRune(c)
			lx.readRune()
		}
	}

	return nil
}

func (lx *Lexer) readDollarExpansion(peekChar rune) error {
	if peekChar != '$' {
		return fmt.Errorf("Invalid usage of readDollarExpansion [bug!]")
	} else {
		lx.yieldNextRuneAsToken(Dollar)
	}

	c, _, err := lx.peekRune()
	if c == '{' {
		err = lx.readDollarExpansionBrace(c)
	} else if c == '(' {
		err = lx.readDollarExpansionParens(c)
	} else if !unicode.IsDigit(c) && IsWordChar(c) {
		// TODO: readAlphaNumeric turns "if" into an If token. but it is legal
		// to have a variable named "if". Ensure this is a NAME token.
		err = lx.readAlphaNumeric(c)
	}
	return err
}

func (lx *Lexer) readDollarExpansionBrace(peekChar rune) error {
	if peekChar != '{' {
		return fmt.Errorf("Invalid usage of readDollarExpansion [bug!}")
	} else {
		lx.yieldNextRuneAsToken(LeftBrace)
	}

	// we must have a name following the left brace
	// TODO: except for ${#VAR} syntax to get string length
	if c, _, err := lx.peekRune(); err != nil {
		return err
	} else if !unicode.IsDigit(c) && IsWordChar(c) {
		err = lx.readAlphaNumeric(c)
		if err != nil {
			return err
		}
	}

	if lx.hasRune('}') {
		lx.yieldNextRuneAsToken(RightBrace)
		return nil
	}

	// TODO: support :-, :=, :?, :+, -, =, ?, + syntax

	// start := lx.location.Copy()
	// var buffer bytes.Buffer
	// for {
	// 	c, _, err := lx.peekRune()
	// 	if err == io.EOF {
	// 		return fmt.Errorf("Unclosed brace expansion")
	// 	} else if err != nil {
	// 		return err
	// 	}

	// 	switch c {
	// 	case '\\':
	// 		lx.readRune()
	// 		// handle an escaped right brace
	// 		if lx.hasRune('}') {
	// 			lx.readRune()
	// 			buffer.WriteRune('}')
	// 		} else {
	// 			buffer.WriteRune('\\')
	// 		}
	// 	case '}':
	// 		text := buffer.String()
	// 		if len(text) > 0 {
	// 			token := &Token{StringSegment, text, start, lx.location.Copy()}
	// 			lx.yieldToken(token)
	// 			start = lx.location.Copy()
	// 			buffer.Reset()
	// 		}

	// 		lx.yieldNextRuneAsToken(RightBrace)
	// 		return nil
	// 	default:
	// 		lx.readRune()
	// 		buffer.WriteRune(c)
	// 	}
	// }

	return nil
}

func (lx *Lexer) yieldNextRuneAsToken(ttype TokenType) {
	start := lx.location.Copy()
	c, _, err := lx.readRune()
	if err == nil {
		text := fmt.Sprintf("%c", c)
		token := &Token{ttype, text, start, lx.location.Copy()}
		lx.peekbuf = append(lx.peekbuf, &peekbufEntry{token, nil})
	}
}

func (lx *Lexer) readDollarExpansionParens(peekChar rune) error {
	return fmt.Errorf("Not implemented")
}

func (lx *Lexer) detectReservedWord(token *Token) (TokenType, bool) {
	var ttype TokenType = Unknown
	switch token.Text {
	case "if":
		ttype = If
	case "then":
		ttype = Then
	case "else":
		ttype = Else
	case "elif":
		ttype = Elif
	case "fi":
		ttype = Fi
	case "do":
		ttype = Do
	case "done":
		ttype = Done
	case "case":
		ttype = Case
	case "esac":
		ttype = Esac
	case "while":
		ttype = While
	case "until":
		ttype = Until
	case "for":
		ttype = For
	case "function":
		ttype = Function
	case "in":
		ttype = In
	}
	return ttype, ttype != Unknown
}

func (lx *Lexer) possiblePunctation(c rune) bool {
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
