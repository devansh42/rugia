package main

import (
	"bytes"
	"fmt"
)

type TokenType uint8
type valBytes []byte
type Token struct {
	tokenType TokenType
	val       valBytes
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %d, Literal: %s", t.tokenType, string(t.val))
}

func (t Token) Equals(tt Token) bool {
	return t.tokenType == tt.tokenType && bytes.Equal(t.val, tt.val)
}

func NewToken(ttype TokenType, val ...byte) Token {
	return Token{tokenType: ttype, val: val}
}

const (
	Plus        TokenType = iota + 1 // +
	Minus                            // -
	Asterisk                         // *
	Slash                            // /
	CommentHash                      // # (for comment)
	Assign                           // =
	Eq                               // ==
	Not                              // !
	NotEq                            // !=
	LT                               // <
	GT                               // >
	LTEQ                             // <=
	GTEQ                             // >=
	RaisePower                       // **
	Mod                              // %
	Number                           // 1234, 12.23
	SemiColon                        // ;
	LeftBrace                        // (
	RightBrace                       // )
)
