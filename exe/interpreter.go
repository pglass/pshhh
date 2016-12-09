package exe

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pglass/pshhh/ast"
	"github.com/pglass/pshhh/lex"
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
	case *ast.Str:
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

func (i *Interpreter) interpretString(node *ast.Str) string {
	var buffer bytes.Buffer
	for _, piece := range node.Pieces {
		switch p := piece.(type) {
		case ast.RawStr:
			buffer.WriteString(string(p))
		case *ast.ParameterExpansion:
			word_val := ""
			if p.Word != nil {
				word_val = i.interpretString(p.Word)
			}
			substitution := i.getParamExpansionSubstitution(p, word_val)
			log.Printf("Evaluated Param Expansion: ${%v} -> %q", p.VarName.Text, substitution)
			// TODO: support :=, :-, etc
			buffer.WriteString(substitution)
		default:
			log.Fatalf("Unhandled StringPiece type %v", p)
		}
	}

	return buffer.String()
}

func (i *Interpreter) getParamExpansionSubstitution(p *ast.ParameterExpansion, word_val string) string {
	key := p.VarName.Text
	param_is_set, param_val := i.FetchEnvVar(key)
	param_is_null := len(param_val) == 0

	if p.Operator == nil {
		return param_val
	}
	switch p.Operator.Type {
	case lex.ColonDash:
		if param_is_set && !param_is_null {
			return param_val
		} else {
			return word_val
		}
	case lex.Dash:
		if param_is_set {
			return param_val
		} else {
			return word_val
		}
	case lex.Plus:
		if param_is_set {
			return word_val
		}
		return ""
	case lex.ColonPlus:
		if param_is_set && !param_is_null {
			return word_val
		}
		return ""
	case lex.Question:
		if param_is_set {
			return param_val
		} else {
			err_msg := fmt.Sprintf("%v: %v", key, word_val)
			i.exit(err_msg, 1)
		}
	case lex.ColonQuestion:
		if param_is_set && !param_is_null {
			return param_val
		} else {
			err_msg := fmt.Sprintf("%v: %v", key, word_val)
			i.exit(err_msg, 1)
		}
	default:
		log.Printf("WARNING: Unhandled param expansion operator %v", p.Operator)
	}
	return ""
}

/* Env stores environment variables as a list of "<key>=<value>" strings. This
 * fetches the <value> portion given the <key>, or returns empty string.
 *
 * This function returns true if the key is set, and false otherwise.
 */
func (i *Interpreter) FetchEnvVar(key string) (bool, string) {
	key = key + "="
	for _, item := range i.Env {
		if strings.HasPrefix(item, key) {
			return true, strings.SplitN(item, "=", 2)[1]
		}
	}
	return false, ""
}

func (i *Interpreter) exit(err_msg string, code int) {
	// todo: print to stderr?
	fmt.Printf("error: %s\n", err_msg)
	os.Exit(1)
}
