package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pglass/pshhh/ast"
	"github.com/pglass/pshhh/exe"
	"github.com/pglass/pshhh/lex"
)

type EnvVars []string

func (e *EnvVars) String() string {
	return fmt.Sprint(*e)
}

func (e *EnvVars) Set(value string) error {
	*e = append(*e, value)
	return nil
}

var (
	filename string
	text     string
	debug    bool
	env_vars EnvVars
)

func init() {
	flag.StringVar(&filename, "f", "", "The filename of a script to run")
	flag.StringVar(&text, "t", "", "Execute the given text")
	flag.BoolVar(&debug, "d", false, "Enable debug mode")
	flag.Var(&env_vars, "e", "Preset environment variables")
}

func main() {
	flag.Parse()

	log.SetFlags(0)
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	var lexer *lex.Lexer = nil
	if filename != "" {
		if b, err := ioutil.ReadFile(filename); err != nil {
			log.Fatal(err)
		} else {
			lexer = lex.NewLexer(string(b))
		}
	} else if text != "" {
		lexer = lex.NewLexer(text)
	}

	if lexer != nil {
		if err := run_single(lexer); err != nil {
			handle_error(err, false)
		}
	} else {
		run_shell()
	}
}

func run_single(lexer *lex.Lexer) error {
	parser := ast.NewParser(lexer)
	root, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("%v\n", err)
	} else if root == nil {
		return fmt.Errorf("Parse failure (got nil node)\n")
	} else {
		interpreter := exe.NewInterpreter()
		interpreter.Env = env_vars

		log.Printf("Environment:")
		for _, item := range interpreter.Env {
			log.Printf("  %v", item)
		}

		if err := interpreter.Interpret(root); err != nil {
			return err
		}
	}
	return nil
}

func run_shell() {
	for {
		fmt.Print("$ ")

		reader := bufio.NewReader(os.Stdin)
		if input, err := reader.ReadString('\n'); err != nil {
			log.Fatalf("error: %v\n", err)
		} else {
			lexer := lex.NewLexer(strings.TrimSpace(input))
			if err := run_single(lexer); err != nil {
				handle_error(err, true)
			}
		}
	}
}

func handle_error(err error, is_shell bool) {
	fmt.Print(err)
	switch e := err.(type) {
	case exe.ExitError:
		if !is_shell {
			os.Exit(e.ExitCode)
		}
	default:
		if !is_shell {
			os.Exit(1)
		}
	}
}
