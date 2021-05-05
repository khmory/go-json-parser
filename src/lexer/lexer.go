package lexer

import (
	"errors"
)

type Lexer struct {
	json     string
	length   int
	position int
}

type LeftCurlyBracketToken struct {
	value string
}

type RightCurlyBracketToken struct {
	value string
}

type LeftSquareBracketToken struct {
	value string
}

type RightSquareBracketToken struct {
	value string
}

type ColonToken struct {
	value string
}

type CommaToken struct {
	value string
}

var (
	charNotFoundError     = errors.New("char is not found")
	invalidCharacterError = errors.New("invalid character")
)

func NewLexer(json string) *Lexer {
	l := new(Lexer)
	l.json = json
	l.length = len(json)
	l.position = 0
	return l
}

func (l Lexer) GetNextToken() (interface{}, error) {
	ch := l.current()
	for {
		if !isSkipCharacter(ch) {
			break
		}

		var err error
		ch, err = l.consume()
		if err == charNotFoundError {
			return "", charNotFoundError
		}
	}

	switch ch {
	case "{":
		return &LeftCurlyBracketToken{value: "{"}, nil
	case "}":
		return &RightCurlyBracketToken{value: "}"}, nil
	case "[":
		return &LeftSquareBracketToken{value: "["}, nil
	case "]":
		return &RightSquareBracketToken{value: "]"}, nil
	case ":":
		return &ColonToken{value: ":"}, nil
	case ",":
		return &CommaToken{value: ","}, nil
	default:
		return "", invalidCharacterError
	}
}

func isSkipCharacter(ch string) bool {
	if ch == " " || ch == "\n" || ch == "\r" || ch == "\t" {
		return true
	}
	return false
}

func (l Lexer) current() string {
	return l.json[l.position : l.position+1]
}

func (l Lexer) consume() (string, error) {
	if l.length <= l.position {
		return "", charNotFoundError
	}
	l.position++
	token := l.current()
	return token, nil
}
