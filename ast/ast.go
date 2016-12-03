package ast

//
//import (
//	"psh/lex"
//)
//
//type Ast struct {
//	Root  Node
//	lexer *lex.Lexer
//}
//
//func NewAst(lexer *lex.Lexer) *Ast {
//	return &Ast{
//		Root:  &RootNode{},
//		lexer: lexer,
//	}
//}
//
//func (a *Ast) Build() error {
//	for {
//		token, err := a.lexer.Next()
//		if token != nil {
//			fmt.Printf("%v", token)
//		}
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
