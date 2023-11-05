package main

import (
	"bytes"
	"testing"
)

func TestAST(t *testing.T) {
	lexer := NewLexer(bytes.NewReader([]byte("(1)+(2)")))

	err := lexer.Analyze()
	if err != nil {
		t.Error("error while doing lexical analysis: ", err)
		t.Fail()
		return
	}
	t.Logf("Token List: %v", lexer.tokenList)
	actual := parseTokens(lexer.tokenList, parseExpr)
	t.Log("Actual: ", actual.String())
	expected := newASTNode(NewToken(Plus, '+'),

		newASTNode(NewToken(Number, '1'), nil, nil, nil),
		newASTNode(NewToken(Number, '2'), nil, nil, nil), nil)

	if !compareAST(expected, actual) {
		t.Error("Expected AST didn't match with Actual AST")
		t.Error("Expected: ", expected.String())
	}
}

func compareAST(x, y *astNode) bool {

	if x == nil && y == nil {
		return true
	}

	if x != nil && y != nil {
		return x.token.Equals(y.token) && compareAST(x.left, y.left) && compareAST(x.right, y.right)
	}
	return false

}
