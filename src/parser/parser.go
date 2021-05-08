package parser

import (
	"errors"
	"fmt"

	"github.com/koheimorii/go-json-parser/src/lexer"
)

var (
	parseError = errors.New("Unparsed tokens detected'")
)

func Parse(json string) (interface{}, error) {
	lex := lexer.NewLexer(json)
	ret, err := ParseValue(lex)
	if err != nil {
		return "", err
	}
	endCheck, _ := lex.GetNextToken()
	if endCheck.TokenType == "EOF" {
		return ret, nil
	}

	fmt.Printf("[debug:%v]", ret)
	return "", parseError
}

func ParseValue(l *lexer.Lexer) (interface{}, error) {
	token, err := l.GetNextToken()
	if err != nil {
		return "", err
	}

	switch token.TokenType {
	case "LeftSquareBracket":
		// todo
	case "LeftCurlyBracket":
		// todo
	case "String", "Number", "True", "False", "Null":
		return token.Value, nil
	}
	fmt.Printf("[debug:%v]", token)
	return "", parseError
}
