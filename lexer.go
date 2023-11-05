package main

import (
	"bufio"
	"errors"
	"io"
)

var illegalToken = errors.New("illegal token")

type Lexer struct {
	reader           *bufio.Reader
	tokenList        []Token
	unProcessedBytes int
	lastErr          error
}

func NewLexer(reader io.Reader) *Lexer {
	lexer := Lexer{
		reader: bufio.NewReader(reader),
	}

	return &lexer
}

func (lexer *Lexer) addT(token Token) {
	lexer.tokenList = append(lexer.tokenList, token)
}

func (lexer *Lexer) eof() bool {
	return lexer.lastErr != nil && lexer.lastErr == io.EOF
}

func (lexer *Lexer) next() byte {
	by, err := lexer.reader.ReadByte()
	if err != nil {
		lexer.lastErr = err
		return 0
	}
	return by
}

func (lexer *Lexer) peek() byte {
	by, err := lexer.reader.Peek(1)
	if err != nil {
		lexer.lastErr = err
		return 0
	}
	return by[0]
}

func (lexer *Lexer) anaEq() {
	nextByte := lexer.peek()
	if !lexer.eof() && nextByte == '=' {
		lexer.next() // consume
		lexer.addT(NewToken(Eq, '=', '='))
		return
	}
	lexer.addT(NewToken(Assign, '='))
}

func (lexer *Lexer) anaAsterisk() {
	nextByte := lexer.peek()
	if !lexer.eof() && nextByte == '*' {
		lexer.next() // consume
		lexer.addT(NewToken(RaisePower, '*', '*'))
		return
	}
	lexer.addT(NewToken(Asterisk, '*'))
}

func (lexer *Lexer) anaNot() {
	nextByte := lexer.peek()
	if !lexer.eof() && nextByte == '=' {
		lexer.next() // consume
		lexer.addT(NewToken(NotEq, '!', '='))
		return
	}
	lexer.addT(NewToken(Not, '!'))
}

func (lexer *Lexer) anaLT() {
	nextByte := lexer.peek()
	if !lexer.eof() && nextByte == '=' {
		lexer.next() //consume
		lexer.addT(NewToken(LTEQ, '<', '='))
		return
	}
	lexer.addT(NewToken(LT, '<'))

}

func (lexer *Lexer) anaGT() {
	nextByte := lexer.peek()
	if !lexer.eof() && nextByte == '=' {
		lexer.next() //consume
		lexer.addT(NewToken(GTEQ, '>', '='))
		return
	}
	lexer.addT(NewToken(GT, '>'))

}

func (lexer *Lexer) anaComment() {
	lexer.addT(NewToken(CommentHash, '#'))
	for {
		nextByte := lexer.next()
		if lexer.eof() || nextByte == '\n' { // Read until next line
			return
		}
	}
}

func (lexer *Lexer) anaDigits(startingDigit byte) error {
	var number = []byte{startingDigit}
	var fractionStarted bool

	for {
		nextByte := lexer.peek()
		if isDigit(nextByte) {
			lexer.next() // consume
			number = append(number, nextByte)
		} else if nextByte == '.' {
			lexer.next() //consume
			if fractionStarted {
				return illegalToken
			}
			number = append(number, nextByte)
			fractionStarted = true

		} else {
			lexer.addT(NewToken(Number, number...))
			return nil
		}
	}

}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isWhiteSpace(b byte) bool {
	return b == '\n' || b == ' ' || b == '\t' || b == '\r'
}

func (lexer *Lexer) Analyze() error {
	for nextByte := lexer.next(); lexer.lastErr == nil; nextByte = lexer.next() {

		switch nextByte {
		case '+':
			lexer.addT(NewToken(Plus, nextByte))
		case '-':
			lexer.addT(NewToken(Minus, nextByte))
		case '%':
			lexer.addT(NewToken(Mod, nextByte))
		case ';':
			lexer.addT(NewToken(SemiColon, nextByte))
		case '#':
			lexer.anaComment()
		case '/':
			lexer.addT(NewToken(Slash, nextByte))
		case '(':
			lexer.addT(NewToken(LeftBrace, '('))
		case ')':
			lexer.addT(NewToken(RightBrace, ')'))
		case '=':
			lexer.anaEq()
		case '*':
			lexer.anaAsterisk()
		case '<':
			lexer.anaLT()
		case '>':
			lexer.anaGT()
		case '!':
			lexer.anaNot()
		default:
			if isDigit(nextByte) {
				err := lexer.anaDigits(nextByte)
				if err != nil { // illegal number
					return err
				}
			}
		}
	}
	if lexer.lastErr == io.EOF {
		return nil
	}
	return lexer.lastErr

}
