//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
)

var (
	_ Expr = IntegerValue(0)
	_ Expr = &binary{}
)

// Expr implements an expression.
type Expr interface {
	Eval() (Value, error)
}

func parseExpr() (Expr, error) {
	return parseLogicalOR()
}

func parseLogicalOR() (Expr, error) {
	return parseLogicalAND()
}

func parseLogicalAND() (Expr, error) {
	return parseBitwiseOR()
}

func parseBitwiseOR() (Expr, error) {
	return parseBitwiseXOR()
}

func parseBitwiseXOR() (Expr, error) {
	return parseBitwiseAND()
}

func parseBitwiseAND() (Expr, error) {
	return parseEquality()
}

func parseEquality() (Expr, error) {
	return parseRelational()
}

func parseRelational() (Expr, error) {
	return parseShift()
}

func parseShift() (Expr, error) {
	return parseAdditive()
}

func parseAdditive() (Expr, error) {
	left, err := parseUnary()
	if err != nil {
		return nil, err
	}
	if !input.HasToken() {
		return left, nil
	}
	t, err := input.GetToken()
	if err != nil {
		return nil, err
	}
	switch t.Type {
	case TAdd, TSub:

	default:
		input.UngetToken(t)
		return left, nil
	}
	right, err := parseUnary()
	if err != nil {
		return nil, err
	}
	return &binary{
		op:    t.Type,
		col:   t.Column,
		left:  left,
		right: right,
	}, nil
}

func parseMultiplicative() (Expr, error) {
	left, err := parseUnary()
	if err != nil {
		return nil, err
	}
	if !input.HasToken() {
		return left, nil
	}
	t, err := input.GetToken()
	if err != nil {
		return nil, err
	}
	switch t.Type {
	case TMult, TDiv, TPercent:

	default:
		input.UngetToken(t)
		return left, nil
	}
	right, err := parseUnary()
	if err != nil {
		return nil, err
	}
	return &binary{
		op:    t.Type,
		col:   t.Column,
		left:  left,
		right: right,
	}, nil
}

func parseUnary() (Expr, error) {
	return parsePostfix()
}

func parsePostfix() (Expr, error) {
	t, err := input.GetToken()
	if err != nil {
		return nil, err
	}
	switch t.Type {
	case TInteger:
		return IntegerValue(t.IntVal), nil

	default:
		input.UngetToken(t)
		return nil, NewError(t.Column, fmt.Errorf("unexpected token '%s'", t))
	}
}

type binary struct {
	op    TokenType
	col   int
	left  Expr
	right Expr
}

func (b binary) String() string {
	return fmt.Sprintf("%s %s %s", b.left, b.op, b.right)
}

func (b binary) Eval() (Value, error) {
	switch l := b.left.(type) {
	case IntegerValue:
		r, ok := b.right.(IntegerValue)
		if !ok {
			return nil,
				NewError(b.col,
					fmt.Errorf("invalid '%s' operation between %T and %T",
						b.op, b.left, b.right))
		}
		var result IntegerValue
		switch b.op {
		case TDiv:
			result = l / r

		case TMult:
			result = l * r

		case TPercent:
			result = l % r

		case TAdd:
			result = l + r

		case TSub:
			result = l - r

		default:
			return nil,
				NewError(b.col,
					fmt.Errorf("unsupport binary operand '%s'", b.op))
		}
		return result, nil

	default:
		return nil,
			NewError(b.col,
				fmt.Errorf("invalid '%s' operation between %T and %T",
					b.op, b.left, b.right))
	}
}
