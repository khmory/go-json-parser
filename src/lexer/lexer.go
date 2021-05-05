package lexer

import (
	"errors"
)

type Lexer struct {
	json     string
	length   int
	position int
}

var charNotFoundError = errors.New("char is not found")

func NewLexer(json string) *Lexer {
	l := new(Lexer)
	l.json = json
	l.length = len(json)
	l.position = 0
	return l
}

func (l Lexer) GetNextToken() (string, error) {
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
	return ch, nil
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
