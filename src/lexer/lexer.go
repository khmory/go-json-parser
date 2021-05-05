package lexer

import (
	"errors"
)

type Lexer struct {
	json     string
	length   int
	position int
}

func NewLexer(json string) *Lexer {
	l := new(Lexer)
	l.json = json
	l.length = len(json)
	l.position = 0
	return l
}

func (l Lexer) Current() string {
	return l.json[l.position : l.position+1]
}

func (l Lexer) Consume() (string, error) {
	if l.length <= l.position {
		return "", errors.New("finished")
	}
	l.position++
	token := l.Current()
	return token, nil
}
