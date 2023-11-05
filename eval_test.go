package main

import (
	"bytes"
	"testing"
)

func TestEval(t *testing.T) {
	lexer := NewLexer(bytes.NewReader([]byte("(2+3)*(6/2)*(4**4**(4-3))")))

	err := lexer.Analyze()
	if err != nil {
		t.Error("error while doing lexical analysis: ", err)
		t.Fail()
		return
	}
	t.Logf("Token List: %v", lexer.tokenList)
	ast := parseTokens(lexer.tokenList, parseExpr)
	t.Log("AST: ", ast.String())
	actualResult, err := eval(ast)
	if err != nil {
		t.Error("error occured while evaluting ast: ", err)
		return
	}
	expectedVal := 3840.0
	if actualResult != expectedVal {
		t.Errorf("Wrong Evaluation Actual: %v, Expected: %v", actualResult, expectedVal)
	}
}
