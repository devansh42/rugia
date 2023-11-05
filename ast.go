package main

import (
	"io"
	"log"
	"strings"
)

type astNode struct {
	left, right *astNode
	token       Token
	depth       int // depth in the ast, helps in stringify
	evalr       evaluator
}

type parserFn func(*parser) *astNode

var nilASTNode *astNode = nil

func newASTNode(token Token, left, right *astNode, evalr evaluator) *astNode {
	return &astNode{left, right, token, 0, evalr}
}

type parser struct {
	tokens         []Token
	processedToken int
	lastErr        error
}

func (ast *astNode) String() string {
	marginStr := ""
	var out strings.Builder
	for leftMargin := ast.depth; leftMargin > 0; leftMargin-- {
		marginStr += "\t"
	}
	out.WriteString("\n")
	out.WriteString(marginStr)
	out.WriteString("Token: ")
	out.WriteString(ast.token.String())
	if ast.left != nilASTNode {
		ast.left.depth = ast.depth + 1
		out.WriteString("\n")
		out.WriteString(marginStr)
		out.WriteString("Left: ")
		out.WriteString(ast.left.String())
	}
	if ast.right != nilASTNode {
		ast.right.depth = ast.depth + 1
		out.WriteString("\n")
		out.WriteString(marginStr)
		out.WriteString("Right: ")
		out.WriteString(ast.right.String())

	}
	return out.String()
}

func (p *parser) next() Token {

	if p.processedToken < len(p.tokens) {
		token := p.tokens[p.processedToken]
		p.processedToken++
		return token
	}
	p.lastErr = io.EOF
	return Token{}
}

func (p *parser) peek() Token {
	if p.processedToken < len(p.tokens) {
		return p.tokens[p.processedToken]
	}
	p.lastErr = io.EOF
	return Token{}
}

func (p *parser) eof() bool {
	defer log.Print("LastError: ", p.lastErr)
	return p.lastErr == io.EOF
}

func parseTokens(tokens []Token, prsFn parserFn) *astNode {
	par := parser{
		tokens: tokens,
	}
	return prsFn(&par)
}
