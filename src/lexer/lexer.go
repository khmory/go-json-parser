package lexer

import (
	"errors"
)

type Lexer struct {
	json     string
	length   int
	position int
}
type Token struct {
	TokenType string
	Value     string
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

func getStringToken(l *Lexer) (string, error) {
	var str string
	for {
		ch := l.consume()
		if ch == "\"" {
			return str, nil
		}
		if ch != "\\" {
			str += ch
			continue
		}

		secoundCh := l.consume()
		switch secoundCh {
		case "\"":
			str += "\""
		case "\\":
			str += "\""
		case "/":
			str += "/"
		case "b":
			str += "\b"
		case "f":
			str += "\f"
		case "n":
			str += "\n"
		case "r":
			str += "\r"
		case "t":
			str += "\t"
		case "u":
			str += getCharacterByCodePoint(l)
		}
	}
}

func getCharacterByCodePoint(l *Lexer) (codePoint string) {
	for i := 0; i < 4; i++ {
		ch := l.consume()
		if ch != "EOF" && ("0" <= ch && ch <= "9") || ("A" <= ch && ch <= "F") || ("a" <= ch && ch <= "f") {
			codePoint += ch
		}
	}
	return
}

func (l *Lexer) GetNextToken() (*Token, error) {
	for {
		ch := l.consume()

		if ch == "EOF" {
			return &Token{TokenType: "EOF", Value: "EOF"}, nil
		}

		if isSkipCharacter(ch) {
			continue
		}

		switch ch {
		case "{":
			return &Token{TokenType: "RightCurlyBracket", Value: "{"}, nil
		case "}":
			return &Token{TokenType: "LeftCurlyBracket", Value: "}"}, nil
		case "[":
			return &Token{TokenType: "LeftSquareBracket", Value: "["}, nil
		case "]":
			return &Token{TokenType: "RightSquareBracket", Value: "]"}, nil
		case ":":
			return &Token{TokenType: "Colon", Value: ":"}, nil
		case ",":
			return &Token{TokenType: "Comma", Value: ","}, nil
		case "\"":
			stringToken, err := getStringToken(l)
			if err != nil {
				return &Token{}, invalidStringError
			}
			return &Token{TokenType: "String", Value: stringToken}, nil
		default:
			return &Token{}, invalidCharacterError
		}
	}
}

func isSkipCharacter(ch string) bool {
	if ch == " " || ch == "\n" || ch == "\r" || ch == "\t" {
		return true
	}
	return false
}

func (l *Lexer) current() string {
	return l.json[l.position : l.position+1]
}

func (l *Lexer) consume() string {
	if l.length <= l.position {
		return "EOF"
	}
	ch := l.current()
	l.position++
	return ch
}
