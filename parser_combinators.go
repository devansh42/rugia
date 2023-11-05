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

// unary -> (-) unary | priamry
func parseUnary(pr *parser) *astNode {
	log.Print("inside unary")
	token := pr.peek()
	if pr.eof() {
		return nil
	}
	if token.tokenType == Minus { // (-) unary
		pr.next() // consuming minus
		// not checking for eof as we have checked while calling peek()
		return newASTNode(token, newASTNode(NewToken(Number, '0'), nilASTNode, nilASTNode, evalNumber), parseUnary(pr), evalInfix)
	}

	return parsePrimary(pr)

}

// factor -> unary ( ( "/" | "*" | "%" ) unary)*
func parseFactor(pr *parser) *astNode {
	log.Print("inside factor")

	leftUnary := parseUnary(pr) // consume first unary
	var parent = leftUnary

	for token := pr.peek(); token.tokenType == Mod || token.tokenType == Asterisk || token.tokenType == Slash; token = pr.peek() {

		if pr.eof() {
			break // we are just breaking the rule as we in the recusrive expansion of (*|/) unary
		}
		pr.next() // consume  * or /
		// not checkig for eof as we have checked while calling peek()
		rightUnary := parseUnary(pr)

		parent = newASTNode(token, parent, rightUnary, evalInfix)
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
		rightUnary := parseFactor(pr)

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
			rightUnary := parseTerm(pr)

			parent = newASTNode(token, parent, rightUnary, evalInfix)
		default:

			return parent
		}

	}
	return parent
}

// expr -> comparison ( ( "==" ) comparison )*
func parseExpr(pr *parser) *astNode {
	log.Print("inside expr")

	leftUnary := parseComp(pr) // consume first unary
	var parent = leftUnary

	for token := pr.peek(); token.tokenType == Eq; token = pr.peek() {
		if pr.eof() {
			break // we are just breaking the rule as we in the recusrive expansion of (==) comparison
		}

		pr.next() // consume ==
		rightUnary := parseComp(pr)

		parent = newASTNode(token, parent, rightUnary, evalInfix)

	}
	return parent

}
