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

func parseExponent(pr *parser) *astNode {
	log.Print("inside exponent")

	leftUnary := parseUnary(pr) // consume first unary
	var parent = leftUnary

	for token := pr.peek(); token.tokenType == RaisePower; token = pr.peek() {

		if pr.eof() {
			break // we are just breaking the rule as we in the recusrive expansion of (**) unary
		}
		pr.next() // consume  **
		// not checkig for eof as we have checked while calling peek()
		rightUnary := parseExponent(pr)

		parent = newASTNode(token, parent, rightUnary, evalInfix) // exponent is right associated operation for us
	}
	return parent
}

// factor -> exponent ( ( "/" | "*" | "%" ) exponent)*
func parseFactor(pr *parser) *astNode {
	log.Print("inside factor")

	leftExponent := parseExponent(pr) // consume first unary
	var parent = leftExponent

	for token := pr.peek(); token.tokenType == Mod || token.tokenType == Asterisk || token.tokenType == Slash; token = pr.peek() {

		if pr.eof() {
			break // we are just breaking the rule as we in the recusrive expansion of (*|/) unary
		}
		pr.next() // consume  * or /
		// not checkig for eof as we have checked while calling peek()
		rightExponent := parseFactor(pr) // Makr

		parent = newASTNode(token, parent, rightExponent, evalInfix)
	}
	return parent
}

// term -> factor (( "+" | "-" ) factor)*
func parseTerm(pr *parser) *astNode {
	log.Print("inside term")
	leftUnary := parseFactor(pr) // consume first unary
	var parent = leftUnary
	log.Print("Peek Token : ", pr.peek())
	for token := pr.peek(); token.tokenType == Plus || token.tokenType == Minus; token = pr.peek() {
		log.Print("Inside Loop")
		if pr.eof() {
			break // we are just breaking the rule as we in the recusrive expansion of (+|-) factor
		}

		pr.next() // consume + or -
		rightUnary := parseTerm(pr)

		parent = newASTNode(token, parent, rightUnary, evalInfix)

	}
	log.Print("Out of Loop")
	return parent

}

// comparison ->  term (( ">" | ">=" | "<" | "<=" ) term)*
func parseComp(pr *parser) *astNode {
	log.Print("inside comp")
	leftUnary := parseTerm(pr) // consume first unary
	var parent = leftUnary
	for {
		token := pr.peek()
		if pr.eof() {
			break // we are just breaking the rule as we in the recusrive expansion of (>|>=|<|<=) term
		}

		switch token.tokenType {
		case LT, GT, LTEQ, GTEQ:
			pr.next() // consume  > or >= or < or <=
			rightUnary := parseComp(pr)

			parent = newASTNode(token, parent, rightUnary, evalInfix)
		default:

			return parent
		}

	}
	return parent
}

// expr -> comparison ( ( "==" | "!=" ) comparison )*
func parseExpr(pr *parser) *astNode {
	log.Print("inside expr")

	leftUnary := parseComp(pr) // consume first unary
	var parent = leftUnary

	for token := pr.peek(); token.tokenType == Eq || token.tokenType == NotEq; token = pr.peek() {
		if pr.eof() {
			break // we are just breaking the rule as we in the recusrive expansion of (== or !=) comparison
		}

		pr.next() // consume == or !=
		rightUnary := parseExpr(pr)

		parent = newASTNode(token, parent, rightUnary, evalInfix)

	}
	return parent

}
