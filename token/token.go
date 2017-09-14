package token

import (
	"fmt"
)

// Token is a token
type Token struct {
	TokenType int
	Lexeme    string
	Literal   interface{}
	Line      int
}

func (t Token) String() string {
	return fmt.Sprintf("%d %s %v", t.TokenType, t.Lexeme, t.Literal)
}
