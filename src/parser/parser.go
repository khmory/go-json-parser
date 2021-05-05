package parser

import (
	"github.com/koheimorii/go-json-parser/src/lexer"
)

func Parse(json string) (parsed string) {
	lex := lexer.NewLexer(json)
	lex.Current()
	return
}
