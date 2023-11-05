package main

import (
	"bytes"
	"testing"
)

func TestLexer(t *testing.T) {
	testMap := map[string][]Token{
		"1.3+2.343": {
			NewToken(Number, []byte("1.3")...),
			NewToken(Plus, '+'),
			NewToken(Number, []byte("2.343")...),
		},
		"(<><=>= ===> =)": {
			NewToken(LeftBrace, '('),
			NewToken(LT, '<'),
			NewToken(GT, '>'),
			NewToken(LTEQ, '<', '='),
			NewToken(GTEQ, '>', '='),
			NewToken(Eq, '=', '='),
			NewToken(Assign, '='),
			NewToken(GT, '>'),
			NewToken(Assign, '='),
			NewToken(RightBrace, ')'),
		},
	}
	for k, v := range testMap {
		if !matchLexialAnalysis(t, k, v) {
			t.Fail()
			return
		}
	}

}

func matchLexialAnalysis(t *testing.T, expr string, expectedTokens []Token) bool {
	lexer := NewLexer(bytes.NewReader([]byte(expr)))
	err := lexer.Analyze()
	if err != nil {
		t.Error("error occured in lexical analysis ", err)
	}
	for i := range expectedTokens {

		if len(expectedTokens) != len(lexer.tokenList) {
			t.Errorf("varing list of tokens, expected: %d, actual: %d", (expectedTokens), (lexer.tokenList))
			return false
		}
		expected := expectedTokens[i]
		actual := lexer.tokenList[i]
		if !expected.Equals(actual) {
			t.Errorf("%s didn't matched with expected %s", actual, expected)
			return false
		}
	}
	return true
}
