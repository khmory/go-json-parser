package parser

import (
	"errors"
	"fmt"

	"github.com/koheimorii/go-json-parser/src/lexer"
)

var (
	parseError = errors.New("Unparsed tokens detected'")
)

const (
	STATE_START = "START"
	STATE_COMMA = "COMMA"
	STATE_VALUE = "VALUE"
)

func Parse(json string) (interface{}, error) {
	lex := lexer.NewLexer(json)
	token, err := lex.GetNextToken()
	if err != nil {
		return nil, err
	}
	ret, err := parseValue(lex, token)
	if err != nil {
		return nil, err
	}
	endCheck, _ := lex.GetNextToken()
	if endCheck.TokenType == "EOF" {
		return ret, nil
	}

	fmt.Printf("[debug:%v]", ret)
	return "", parseError
}

func parseValue(l *lexer.Lexer, token *lexer.Token) (interface{}, error) {
	switch token.TokenType {
	case "LeftSquareBracket":
		mapValue, err := parseArray(l)
		if err != nil {
			return nil, err
		}
		return mapValue, nil
	case "LeftCurlyBracket":
		// todo
	case "String", "Number", "True", "False", "Null":
		return token.Value, nil
	}
	fmt.Printf("[debug:%v]", token)
	return "", parseError
}

func parseArray(l *lexer.Lexer) ([]interface{}, error) {
	var array []interface{}
	state := STATE_START

	for {
		token, err := l.GetNextToken()
		if err != nil {
			return nil, err
		}
		if token.TokenType == "EOF" {
			break
		}

		switch state {
		case STATE_START:
			if token.TokenType == "RightSquareBracket" {
				return array, nil
			}
			value, err := parseValue(l, token)
			if err != nil {
				return nil, err
			}
			array = append(array, value)
			state = STATE_VALUE
		case STATE_VALUE:
			if token.TokenType == "RightSquareBracket" {
				return array, nil
			}
			if token.TokenType == "Comma" {
				state = STATE_COMMA
				break
			}
			fmt.Println("invalid token in value state")
			return nil, parseError
		case STATE_COMMA:
			value, err := parseValue(l, token)
			if err != nil {
				return nil, err
			}
			array = append(array, value)
			state = STATE_VALUE
		default:
			fmt.Println("invalid state in array parser")
			return nil, parseError
		}
	}
	return nil, parseError
}
