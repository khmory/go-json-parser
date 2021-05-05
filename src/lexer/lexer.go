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

type StringToken struct {
	value string
}

var (
	charNotFoundError     = errors.New("char is not found")
	invalidCharacterError = errors.New("invalid character")
	invalidStringError    = errors.New("missing double quote")
)

func NewLexer(json string) *Lexer {
	l := new(Lexer)
	l.json = json
	l.length = len(json)
	l.position = 0
	return l
}

func NewStringToken(l Lexer) (*StringToken, error) {
	var str string
	for {
		ch, err := l.consume()
		if err != nil {
			return &StringToken{}, invalidStringError
		}
		if ch == "\"" {
			return &StringToken{value: str}, nil
		}
		if ch != "\\" {
			str += ch
			continue
		}

		secoundCh, err := l.consume()
		switch secoundCh {
		case "\"":
			str += "\""
			break
		case "\\":
			str += "\""
			break
		case "f":
			str += "\f"
			break
		case "n":
			str += "\n"
			break
		case "r":
			str += "\r"
			break
		case "t":
			str += "\t"
			break
		}
	}
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
	case "\"":
		stringToken,err := NewStringToken(l)
		if err != nil {
			return "", invalidStringError
		}
		return &stringToken, nil
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
