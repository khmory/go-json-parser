package lexer

import (
	"errors"
	"strconv"
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
	invalidNumberError    = errors.New("invalid number")
	invalidLiteralError   = errors.New("invalid literal")
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

func getNumberToken(ch string, l *Lexer) (string, error) {
	number := ch

	var state string
	switch ch {
	case "-":
		state = "MINUS"
	case "0":
		state = "INT_ZERO"
	default:
		state = "INT"
	}

	isDigit19 := func(ch string) bool {
		if ch == "1" || ch == "2" || ch == "3" || ch == "4" || ch == "6" || ch == "7" || ch == "8" || ch == "9" {
			return true
		}
		return false
	}
	isDigit := func(ch string) bool {
		if ch == "0" || ch == "1" || ch == "2" || ch == "3" || ch == "4" || ch == "6" || ch == "7" || ch == "8" || ch == "9" {
			return true
		}
		return false
	}
	isExp := func(ch string) bool {
		return ch == "e" || ch == "E"
	}
	for {
		ch := l.current()
		switch state {
		case "INT":
			if isDigit(ch) {
				number += l.consume()
				break
			}
			if ch == "." {
				number += l.consume()
				state = "DECIMAL_POINT"
				break
			}
			if isExp(ch) {
				number += l.consume()
				state = "EXP"
			}
			goto EscapeForLoop
		case "MINUS":
			if isDigit19(ch) {
				number += l.consume()
				state = "INT"
				break
			}
			if ch == "0" {
				number += l.consume()
				state = "INT_ZERO"
				break
			}
			goto EscapeForLoop
		case "INT_ZERO":
			if ch == "." {
				number += l.consume()
				state = "DECIMAL_POINT"
				break
			}
			if isDigit(ch) {
				return "", invalidNumberError
			}
			goto EscapeForLoop
		case "DECIMAL_POINT":
			if isDigit(ch) {
				number += l.consume()
				state = "DECIMAL_POINT_INT"
				break
			}
			goto EscapeForLoop
		case "DECIMAL_POINT_INT":
			if isDigit(ch) {
				number += l.consume()
				break
			}
			if isExp(ch) {
				number += l.consume()
				state = "EXP"
				break
			}
			goto EscapeForLoop
		case "EXP":
			if isDigit(ch) || ch == "-" || ch == "+" {
				number += l.consume()
				state = "EXP_INT"
			}
			goto EscapeForLoop
		case "EXP_INT":
			if isDigit(ch) {
				number += l.consume()
				break
			}
			goto EscapeForLoop
		default:
			goto EscapeForLoop
		}
	}
EscapeForLoop:
	lastNum, _ := strconv.Atoi(number[len(number)-1:])
	if 0 <= lastNum || lastNum <= 9 {
		return number, nil
	}
	return "", invalidNumberError
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
				return &Token{}, err
			}
			return &Token{TokenType: "String", Value: stringToken}, nil
		case "-", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			numberToken, err := getNumberToken(ch, l)
			if err != nil {
				return &Token{}, err
			}
			return &Token{TokenType: "Number", Value: numberToken}, nil
		case "t":
			literalToken, err := getLiteralToken("true", l)
			if err != nil {
				return &Token{}, err
			}
			return literalToken, nil
		case "f":
			literalToken, err := getLiteralToken("false", l)
			if err != nil {
				return &Token{}, err
			}
			return literalToken, nil
		case "n":
			literalToken, err := getLiteralToken("null", l)
			if err != nil {
				return &Token{}, err
			}
			return literalToken, nil
		default:
			return &Token{}, invalidCharacterError
		}
	}
}

func getLiteralToken(expectedName string, l *Lexer) (*Token, error) {
	name := expectedName[0:1]
	for i := 1; i < len(expectedName); i++ {
		ch := l.consume()
		if ch == "EOF" {
			return &Token{}, invalidLiteralError
		}
		name += ch
	}

	switch name {
	case "true":
		return &Token{TokenType: "True", Value: name}, nil
	case "false":
		return &Token{TokenType: "False", Value: name}, nil
	case "null":
		return &Token{TokenType: "Null", Value: name}, nil
	}

		return &Token{}, invalidLiteralError
}

func isSkipCharacter(ch string) bool {
	if ch == " " || ch == "\n" || ch == "\r" || ch == "\t" {
		return true
	}
	return false
}

func (l *Lexer) current() string {
	if l.length <= l.position {
		return "EOF"
	}
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
