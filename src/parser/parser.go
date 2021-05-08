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
	token, err := lex.GetNextToken()
	if err != nil {
		return nil, err
	}
	value, err := parseValue(lex, token)
	if err != nil {
		return nil, err
	}
	endCheck, _ := lex.GetNextToken()
	if endCheck.TokenType == "EOF" {
		return value, nil
	}

	fmt.Printf("[debug:%v]", value)
	return "", parseError
}

func parseValue(l *lexer.Lexer, token *lexer.Token) (interface{}, error) {
	switch token.TokenType {
	case "LeftSquareBracket":
		arrayValue, err := parseArray(l)
		if err != nil {
			return nil, err
		}
		return arrayValue, nil
	case "LeftCurlyBracket":
		mapValue, err := parseMap(l)
		if err != nil {
			return nil, err
		}
		return mapValue, nil
	case "String", "Number", "True", "False", "Null":
		return token.Value, nil
	}
	fmt.Printf("[debug:%v]", token)
	return "", parseError
}

func parseArray(l *lexer.Lexer) ([]interface{}, error) {

	const (
		STATE_START = "START"
		STATE_COMMA = "COMMA"
		STATE_VALUE = "VALUE"
	)

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

func parseMap(l *lexer.Lexer) (map[interface{}]interface{}, error) {
	const (
		STATE_START = "START"
		STATE_KEY   = "KEY"
		STATE_COLON = "COLON"
		STATE_COMMA = "COMMA"
		STATE_VALUE = "VALUE"
	)

	mapValue := map[interface{}]interface{}{}
	var key interface{}
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
			if token.TokenType == "RightCurlyBracket" {
				return mapValue, nil
			}
			if token.TokenType == "String" {
				key = token.Value
				state = STATE_KEY
				break
			}
			fmt.Printf("invalid token: %s \n", token.TokenType)
			return nil, parseError
		case STATE_KEY:
			if token.TokenType == "Colon" {
				state = STATE_COLON
				break
			}
			fmt.Printf("invalid token: %s \n", token.TokenType)
		case STATE_COLON:
			value, err := parseValue(l, token)
			if err != nil {
				return nil, err
			}
			mapValue[key] = value
			state = STATE_VALUE
		case STATE_VALUE:
			if token.TokenType == "RightCurlyBracket" {
				return mapValue, nil
			}
			if token.TokenType == "Comma" {
				state = STATE_COMMA
				break
			}
			fmt.Printf("invalid token: %s \n", token.TokenType)
		case STATE_COMMA:
			if token.TokenType == "String" {
				key = token.Value
				state = STATE_KEY
				break
			}
			fmt.Printf("invalid token: %s \n", token.TokenType)
		default:
			fmt.Println("invalid state in map parser")
			return nil, parseError
		}
	}
	return nil, parseError
}
