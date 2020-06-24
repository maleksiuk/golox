package tokens

import "fmt"

// Token represents what the scanner will produce
type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   interface{}
	Line      int
}

func (token Token) String() string {
	return fmt.Sprintf("%v %v %v", token.TokenType, token.Lexeme, token.Literal)
}
