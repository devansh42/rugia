package main

import (
	"errors"
	"log"
)

// primary -> NUMBER | "(" expr ")" TODO
func parsePrimary(pr *parser) *astNode {
	log.Print("inside primary")
	token := pr.next()
	if pr.eof() {
		return nil
	}
	switch token.tokenType {
	case Number: // NUMBER
		return newASTNode(token, nilASTNode, nilASTNode, evalNumber)
	case LeftBrace: // "(" expr ")"
		ast := parseExpr(pr)
		pr.next()     // consuming ")"
		if pr.eof() { // Because we found eof while expecting ) for consumption
			pr.lastErr = errors.New(") expected")
			return nil
		}
		return ast
	}
	return nil
}

// unary -> (- | !) unary | priamry
func parseUnary(pr *parser) *astNode {
	log.Print("inside unary")
	token := pr.peek()
	if pr.eof() {
		return nil
	}
	if token.tokenType == Minus || token.tokenType == Not {
		pr.next() // consume minus or not
		return newASTNode(token, nilASTNode, parseUnary(pr), evalUnary)
	}

	return parsePrimary(pr)

}

// exponent -> unary ( ("**") unary )*
func parseExponent(pr *parser) *astNode {
	return parseCombinator(pr, evalInfix, parseUnary, RaisePower)
}

// factor -> exponent ( ( "/" | "*" | "%" ) exponent)*
func parseFactor(pr *parser) *astNode {
	return parseCombinator(pr, evalInfix, parseExponent, Mod, Asterisk, Slash)
}

// term -> factor (( "+" | "-" ) factor)*
func parseTerm(pr *parser) *astNode {
	return parseCombinator(pr, evalInfix, parseFactor, Plus, Minus)
}

// comparison ->  term (( ">" | ">=" | "<" | "<=" ) term)*
func parseComp(pr *parser) *astNode {
	return parseCombinator(pr, evalInfix, parseTerm, LT, GT, LTEQ, GTEQ)
}

// expr -> comparison ( ( "==" | "!=" ) comparison )*
func parseExpr(pr *parser) *astNode {
	return parseCombinator(pr, evalInfix, parseComp, Eq, NotEq)
}

func parseCombinator(pr *parser, evalr evaluator, nonterminalParser parserFn, matchableTokens ...TokenType) *astNode {
	var tokenMap = make(map[TokenType]struct{}, len(matchableTokens))
	for _, token := range matchableTokens {
		tokenMap[token] = struct{}{}
	}
	var fn parserFn
	fn = func(pr *parser) *astNode {
		parent := nonterminalParser(pr) // parsing first part of the grammar e.g. in grammar E -> F ( operator E )* parses F
		token := pr.peek()
		_, ok := tokenMap[token.tokenType]
		if !pr.eof() && ok {
			pr.next() // consume token
			return newASTNode(token, parent, fn(pr), evalr)
		}
		return parent
	}
	return fn(pr)
}
