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

func (lexer *Lexer) retriveByte() byte {
	bytes, err := lexer.reader.Peek(lexer.unProcessedBytes)
	if err != nil {
		lexer.lastErr = err
		return 0
	}
	if lexer.unProcessedBytes > 0 {
		return bytes[lexer.unProcessedBytes-1]
	}
	return 0
}

func (lexer *Lexer) eof() bool {
	return lexer.lastErr != nil && lexer.lastErr == io.EOF
}

func (lexer *Lexer) nextByte() byte {
	lexer.unProcessedBytes++
	return lexer.retriveByte()
}
func (lexer *Lexer) prevByte() byte {
	lexer.unProcessedBytes--
	return lexer.retriveByte()
}

func (lexer *Lexer) drain() {
	_, lexer.lastErr = lexer.reader.Discard(lexer.unProcessedBytes)
	lexer.unProcessedBytes = 0
}

func (lexer *Lexer) anaEq() {
	nextByte := lexer.nextByte()
	if !lexer.eof() && nextByte == '=' {
		lexer.addT(NewToken(Eq, '=', '='))
		return
	}
	lexer.addT(NewToken(Assign, '='))
	lexer.prevByte()
}

func (lexer *Lexer) anaAsterisk() {
	nextByte := lexer.nextByte()
	if !lexer.eof() && nextByte == '*' {
		lexer.addT(NewToken(RaisePower, '*', '*'))
		return
	}
	lexer.addT(NewToken(Asterisk, '*'))
	lexer.prevByte()
}

func (lexer *Lexer) anaNot() {
	nextByte := lexer.nextByte()
	if !lexer.eof() && nextByte == '=' {
		lexer.addT(NewToken(NotEq, '!', '='))
		return
	}
	lexer.addT(NewToken(Not, '!'))
	lexer.prevByte()
}

func (lexer *Lexer) anaLT() {
	nextByte := lexer.nextByte()
	if !lexer.eof() && nextByte == '=' {
		lexer.addT(NewToken(LTEQ, '<', '='))
		return
	}
	lexer.addT(NewToken(LT, '<'))
	lexer.prevByte()

}

func (lexer *Lexer) anaGT() {
	nextByte := lexer.nextByte()
	if !lexer.eof() && nextByte == '=' {
		lexer.addT(NewToken(GTEQ, '>', '='))
		return
	}
	lexer.addT(NewToken(GT, '>'))
	lexer.prevByte()

}

func (lexer *Lexer) anaComment() {
	lexer.addT(NewToken(CommentHash, '#'))
	for {
		nextByte := lexer.nextByte()
		if lexer.eof() || nextByte == '\n' { // Read until next line
			return
		}
	}
}

func (lexer *Lexer) anaDigits(startingDigit byte) error {
	var number = []byte{startingDigit}
	var fractionStarted bool

	for {
		nextByte := lexer.nextByte()
		if isDigit(nextByte) {
			number = append(number, nextByte)
		} else if nextByte == '.' {
			if fractionStarted {
				return illegalToken
			}
			number = append(number, nextByte)
			fractionStarted = true

		} else {
			lexer.addT(NewToken(Number, number...))
			lexer.prevByte()
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
	for ; lexer.lastErr == nil; lexer.drain() {
		nextByte := lexer.nextByte()

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
