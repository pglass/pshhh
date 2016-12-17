package exe

import (
	"bytes"
	"fmt"
	"log"
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

func (i *Interpreter) Interpret(node ast.Node) error {
	switch n := node.(type) {
	case *ast.GenericNode:
		return i.interpretGenericNode(n)
	case *ast.CommandList:
		return i.interpretCommandList(n)
	case *ast.Str:
		// a string could be a param expansion that resolves to a program. for example, if you do
		//
		//     $ export FOO='echo wumbo'
		//	   $ $FOO
		//	   wumbo
		//
		// TODO: figure out how this should work. In bash, there seem to be some limitations around
		// what how "programs stored in environment variables" get executed. Are these limited to
		// single commands? Are some limitations for security reasons?
		//
		// For example,
		//
		//     bash> FOO='echo wumbo; echo mini'; $FOO
		//     wumbo; echo mini
		//
		//     bash> export Y=thisisy
		//     bash> FOO='echo $Y'; $FOO
		//     $Y
		//
		if text, err := i.interpretString(n); err != nil {
			return err
		} else {
			lexer := lex.NewLexer(text)
			parser := ast.NewParser(lexer)
			if root, err := parser.Parse(); err != nil {
				return err
			} else {
				return i.Interpret(root)
			}
		}
	}
	return fmt.Errorf("ERROR: Unhandled node %v\n", node)
}

func (i *Interpreter) interpretGenericNode(node *ast.GenericNode) error {
	log.Printf("Interpret GenericNode: %v", node)
	for _, child := range node.Children {
		if err := i.Interpret(child); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) interpretCommandList(node *ast.CommandList) error {
	log.Printf("Interpret CommandList: %v", node)
	for j, command := range node.Commands {
		var proc *PshProc
		var err error

		switch n := command.(type) {
		case *ast.SimpleCommand:
			if proc, err = i.interpretSimpleCommand(n); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unhandled command in CommandList: %v", n)
		}

		// a command is backgrounded if followed by '&'
		// todo: the parser should make this easier for us
		proc.IsBackground = j < len(node.Separators) &&
			node.Separators[j].Type == lex.Ampersand

		if _, err := proc.ForkExec(); err != nil {
			return fmt.Errorf("ERROR: failed to run %v: %v\n", proc.Args, err)
		}
	}
	return nil
}

func (i *Interpreter) interpretSimpleCommand(node *ast.SimpleCommand) (*PshProc, error) {
	log.Printf("Interpret SimpleCommand: %v", node)

	args := []string{}
	for _, word := range node.Words {
		if text, err := i.interpretString(word); err != nil {
			return nil, err
		} else {
			args = append(args, text)
		}
	}

	if proc, err := NewPshProc(args, i.Env); err != nil {
		return nil, err
	} else {
		return proc, nil
	}
}

func (i *Interpreter) interpretString(node *ast.Str) (string, error) {
	var buffer bytes.Buffer
	for _, piece := range node.Pieces {
		switch p := piece.(type) {
		case ast.RawStr:
			buffer.WriteString(string(p))
		case *ast.ParameterExpansion:
			var word_val string
			var err error
			var sub string

			if p.Word != nil {
				if word_val, err = i.interpretString(p.Word); err != nil {
					return "", err
				}
			}

			if sub, err = i.resolveParamExpansion(p, word_val); err != nil {
				return "", err
			} else {
				log.Printf("Evaluated Param Expansion: ${%v} -> %q", p.VarName.Text, sub)
				buffer.WriteString(sub)
			}
		default:
			return "", fmt.Errorf("Unhandled StringPiece type %v", p)
		}
	}

	return buffer.String(), nil
}

func (i *Interpreter) resolveParamExpansion(p *ast.ParameterExpansion, word_val string) (string, error) {
	key := p.VarName.Text
	param_is_set, param_val := i.FetchEnvVar(key)
	param_is_null := len(param_val) == 0

	if p.Operator == nil {
		return param_val, nil
	}
	switch p.Operator.Type {
	case lex.ColonDash:
		if param_is_set && !param_is_null {
			return param_val, nil
		} else {
			return word_val, nil
		}
	case lex.Dash:
		if param_is_set {
			return param_val, nil
		} else {
			return word_val, nil
		}
	case lex.Plus:
		if param_is_set {
			return word_val, nil
		}
		return "", nil
	case lex.ColonPlus:
		if param_is_set && !param_is_null {
			return word_val, nil
		}
		return "", nil
	case lex.Question:
		if param_is_set {
			return param_val, nil
		} else {
			err_msg := fmt.Sprintf("%v: %v", key, word_val)
			return "", i.exit(err_msg, 1)
		}
	case lex.ColonQuestion:
		if param_is_set && !param_is_null {
			return param_val, nil
		} else {
			err_msg := fmt.Sprintf("%v: %v", key, word_val)
			return "", i.exit(err_msg, 1)
		}
	}
	return "", fmt.Errorf("ERROR: Unhandled param expansion operator %v", p.Operator)
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

func (i *Interpreter) exit(err_msg string, code int) error {
	// todo: print to stderr?
	return ExitError{
		error:    fmt.Errorf("error: %s\n", err_msg),
		ExitCode: 1,
	}
}
