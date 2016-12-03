package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"psh/ast"
	"psh/exe"
	"psh/lex"
	"strings"
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

	var lexer *lex.Lexer = nil
	if filename != "" {
		if f, err := os.Open(filename); err != nil {
			log.Fatal(err)
		} else {
			lexer = lex.NewLexer(bufio.NewReader(f))
		}
	} else if text != "" {
		lexer = lex.NewLexer(strings.NewReader(text))
	}

	if lexer == nil {
		flag.PrintDefaults()
		log.Fatal("filename or text required")
	}

	log.SetFlags(0)
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	parser := ast.NewParser(lexer)

	root, err := parser.Parse()
	if err != nil {
		fmt.Printf("%v\n", err)
	} else if root == nil {
		fmt.Printf("Parse failure (got nil node)\n")
	} else {
		interpreter := exe.NewInterpreter()
		interpreter.Env = env_vars

		log.Printf("Environment:")
		for _, item := range interpreter.Env {
			log.Printf("  %v", item)
		}

		interpreter.Interpret(root)
	}
}
