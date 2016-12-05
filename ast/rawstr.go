package ast

import (
	"fmt"
)

type RawStr string

func (s RawStr) IsStrPiece() {}

func (s RawStr) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, string(s))
}
