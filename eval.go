package main

import (
	"errors"
	"strconv"
)

type evaluator func(*astNode) (any, error)

func eval(ast *astNode) (any, error) {
	return ast.evalr(ast)
}

func evalNumber(ast *astNode) (any, error) {
	return strconv.ParseFloat(string(ast.token.val), 64)
}

func evalInfix(ast *astNode) (any, error) {
	var left, right any
	left, err := eval(ast.left)
	if err != nil {
		return nil, err
	}
	right, err = eval(ast.right)
	if err != nil {
		return 0, nil
	}
	leftF := left.(float64)
	rightF := right.(float64)

	switch ast.token.tokenType {
	case Minus:
		return leftF - rightF, nil
	case Plus:
		return leftF + rightF, nil
	case Asterisk:
		return leftF * rightF, nil
	case Slash:
		if rightF == 0 {
			return 0, errors.New("divided by zero")
		}
		return leftF / rightF, nil
	case Mod:
		return float64(int64(leftF) % int64(rightF)), nil

	case LT:
		return leftF < rightF, nil
	case GT:
		return leftF > rightF, nil
	case LTEQ:
		return leftF <= rightF, nil
	case GTEQ:
		return leftF >= rightF, nil
	case Eq:
		return leftF == rightF, nil

	}
	return 0, nil // Need to take care of this
}
