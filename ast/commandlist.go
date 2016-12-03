package ast

import (
	"fmt"
	"psh/lex"
)

/* A list of commands joined by either ';' or '&' */
type CommandList struct {
	Commands   []Command
	Separators []*lex.Token
}

func NewCommandList() *CommandList {
	return &CommandList{
		Commands:   []Command{},
		Separators: []*lex.Token{},
	}
}

func (c *CommandList) Format(f fmt.State, _ rune) {
	fmt.Fprintf(f, "CommandList[")
	if len(c.Commands) > 0 {
		fmt.Fprintf(f, "%v", c.Commands[0])
		for i, command := range c.Commands[1:] {
			if i < len(c.Separators) {
				fmt.Fprintf(f, " %v ", c.Separators[i].Text)
			}
			fmt.Fprintf(f, "%v", command)
		}
	}
	fmt.Fprintf(f, "]")
}

/* Returns an error on failing to parse at least one command.
 * This does not consume newlines. */
func (c *CommandList) Parse(parser *Parser) error {
	for {
		parser.ConsumeWhile(lex.Space)

		// parse the command
		if command, err := c.parseCommand(parser); err != nil {
			return err
		} else if command == nil && len(c.Commands) != 0 {
			break
		} else if command != nil {
			c.Commands = append(c.Commands, command)
		} else {
			return fmt.Errorf("Expected command")
		}

		parser.ConsumeWhile(lex.Space)

		if tok, err := parser.ConsumeAny(lex.Semi, lex.Ampersand); err != nil {
			break
		} else {
			c.Separators = append(c.Separators, tok)
		}
	}
	return nil
}

func (c *CommandList) parseCommand(parser *Parser) (Command, error) {
	tok, err := parser.Lexer.Peek()
	if err != nil {
		return nil, err
	} else if tok == nil {
		return nil, nil
	}

	var command Command = nil
	switch tok.Type {
	case lex.For:
		command = NewForClause()
	case lex.Word, lex.Name, lex.Number:
		command = NewSimpleCommand()
	default:
		return nil, nil
	}

	if err := command.Parse(parser); err != nil {
		return nil, err
	}
	return command, nil
}
