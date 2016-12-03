package exe

import (
	"bytes"
	"fmt"
	"log"
	"psh/ast"
	"psh/lex"
	"strings"
)

type Interpreter struct {
	Debug bool
	Env   []string
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		Debug: false,
	}
}

func (i *Interpreter) Interpret(node ast.Node) {
	// log the environment
	switch n := node.(type) {
	case *ast.GenericNode:
		i.interpretGenericNode(n)
	case *ast.CommandList:
		i.interpretCommandList(n)
	case *ast.String:
		// TODO: a string could be a param expansion that resolves to a command
		// name, e.g. if you do `export FOO=echo; "$FOO"`
		i.interpretString(n)
	default:
		fmt.Printf("ERROR: Unhandled node %v\n", n)
	}
}

func (i *Interpreter) interpretGenericNode(node *ast.GenericNode) {
	log.Printf("Interpret GenericNode: %v", node)
	for _, child := range node.Children {
		i.Interpret(child)
	}
}

func (i *Interpreter) interpretCommandList(node *ast.CommandList) {
	log.Printf("Interpret CommandList: %v", node)
	for j, command := range node.Commands {
		var proc *PshProc

		switch n := command.(type) {
		case *ast.SimpleCommand:
			proc = i.interpretSimpleCommand(n)
		default:
			log.Printf("Unhandled command in CommandList: %v", n)
			return
		}

		// a command is backgrounded if followed by '&'
		// todo: the parser should make this easier for us
		proc.IsBackground = j < len(node.Separators) &&
			node.Separators[j].Type == lex.Ampersand

		if _, err := proc.ForkExec(); err != nil {
			fmt.Printf("ERROR: failed to run %v: %v\n", proc.Args, err)
		}
	}
}

func (i *Interpreter) interpretSimpleCommand(node *ast.SimpleCommand) *PshProc {
	log.Printf("Interpret SimpleCommand: %v", node)

	args := []string{}
	for _, word := range node.Words {
		text := i.interpretString(word)
		args = append(args, text)
	}

	proc, err := NewPshProc(args, i.Env)
	if err != nil {
		log.Fatal(err)
	}
	return proc
}

func (i *Interpreter) interpretString(node *ast.String) string {
	var buffer bytes.Buffer
	for _, piece := range node.Pieces {
		switch p := piece.(type) {
		case *ast.StringSegment:
			buffer.WriteString(p.Text)
		case *ast.ParameterExpansion:
			// TODO: support :=, :-, etc
			key := p.VarName.Text
			val := i.FetchEnvVar(key)
			log.Printf("Evaluated Param Expansion: ${%v} -> %q", key, val)
			buffer.WriteString(val)
		default:
			log.Fatal("Unhandled StringPiece type %v", p)
		}
	}

	return buffer.String()
}

/* c.ProcAttr.Env stores environment variables as a list of "<key>=<value>"
 * strings. This fetches the <value> portion given the <key>, or returns
 * empty string.
 */
func (i *Interpreter) FetchEnvVar(key string) string {
	key = key + "="
	for _, item := range i.Env {
		if strings.HasPrefix(item, key) {
			return strings.SplitN(item, "=", 2)[1]
		}
	}
	return ""
}
