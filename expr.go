//
// Copyright (c) 2020-2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
)

var (
	_ Expr = BoolValue(false)
	_ Expr = Int8Value(0)
	_ Expr = Int16Value(0)
	_ Expr = Int32Value(0)
	_ Expr = Int64Value(0)
	_ Expr = Float64Value(0)
	_ Expr = &binary{}
	_ Expr = &unary{}
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
	left, err := parseAdditive()
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
	case TLeftShift, TRightShift:

	default:
		input.UngetToken(t)
		return left, nil
	}
	right, err := parseAdditive()
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

func parseAdditive() (Expr, error) {
	left, err := parseMultiplicative()
	if err != nil {
		return nil, err
	}
	for {
		if !input.HasToken() {
			return left, nil
		}
		t, err := input.GetToken()
		if err != nil {
			return nil, err
		}
		switch t.Type {
		case '+', '-':

		default:
			input.UngetToken(t)
			return left, nil
		}
		right, err := parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left = &binary{
			op:    t.Type,
			col:   t.Column,
			left:  left,
			right: right,
		}
	}
}

func parseMultiplicative() (Expr, error) {
	left, err := parseUnary()
	if err != nil {
		return nil, err
	}
	for {
		if !input.HasToken() {
			return left, nil
		}
		t, err := input.GetToken()
		if err != nil {
			return nil, err
		}
		switch t.Type {
		case '*', '/', '%':

		default:
			input.UngetToken(t)
			return left, nil
		}
		right, err := parseUnary()
		if err != nil {
			return nil, err
		}
		left = &binary{
			op:    t.Type,
			col:   t.Column,
			left:  left,
			right: right,
		}
	}
}

func parseUnary() (Expr, error) {
	t, err := input.GetToken()
	if err != nil {
		return nil, err
	}
	switch t.Type {
	case '-':
		expr, err := parsePostfix()
		if err != nil {
			return nil, err
		}
		return &unary{
			op:    t.Type,
			col:   t.Column,
			value: expr,
		}, nil

	default:
		input.UngetToken(t)
		return parsePostfix()
	}
}

func parsePostfix() (Expr, error) {
	t, err := input.GetToken()
	if err != nil {
		return nil, err
	}
	switch t.Type {
	case '(':
		expr, err := parseExpr()
		if err != nil {
			return nil, err
		}
		t, err = input.GetToken()
		if err != nil {
			return nil, err
		}
		if t.Type != ')' {
			return nil,
				NewError(t.Column, fmt.Errorf("unexpected token '%s'", t))
		}
		return expr, nil

	case TInteger:
		return t.IntVal, nil

	case TFloat:
		return t.FloatVal, nil

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
	v1, err := b.left.Eval()
	if err != nil {
		return nil, err
	}
	v2, err := b.right.Eval()
	if err != nil {
		return nil, err
	}
	t, err := ConversionType(v1, v2)
	if err != nil {
		return nil, err
	}

	switch t {
	case TypeInt8:
		i1, err := ValueInt8(v1)
		if err != nil {
			return nil, err
		}
		i2, err := ValueInt8(v2)
		if err != nil {
			return nil, err
		}
		var result int8
		switch b.op {
		case '/':
			result = i1 / i2
		case '*':
			result = i1 * i2
		case '%':
			result = i1 % i2
		case '+':
			result = i1 + i2
		case '-':
			result = i1 - i2
		case TLeftShift:
			result = i1 << i2
		case TRightShift:
			result = i1 >> i2
		default:
			return nil,
				NewError(b.col, fmt.Errorf("unsupport binary operand '%s'",
					b.op))
		}
		return Int8Value(result), nil

	case TypeInt16:
		i1, err := ValueInt16(v1)
		if err != nil {
			return nil, err
		}
		i2, err := ValueInt16(v2)
		if err != nil {
			return nil, err
		}
		var result int16
		switch b.op {
		case '/':
			result = i1 / i2
		case '*':
			result = i1 * i2
		case '%':
			result = i1 % i2
		case '+':
			result = i1 + i2
		case '-':
			result = i1 - i2
		case TLeftShift:
			result = i1 << i2
		case TRightShift:
			result = i1 >> i2
		default:
			return nil,
				NewError(b.col, fmt.Errorf("unsupport binary operand '%s'",
					b.op))
		}
		return Int16Value(result), nil

	case TypeInt32:
		i1, err := ValueInt32(v1)
		if err != nil {
			return nil, err
		}
		i2, err := ValueInt32(v2)
		if err != nil {
			return nil, err
		}
		var result int32
		switch b.op {
		case '/':
			result = i1 / i2
		case '*':
			result = i1 * i2
		case '%':
			result = i1 % i2
		case '+':
			result = i1 + i2
		case '-':
			result = i1 - i2
		case TLeftShift:
			result = i1 << i2
		case TRightShift:
			result = i1 >> i2
		default:
			return nil,
				NewError(b.col, fmt.Errorf("unsupport binary operand '%s'",
					b.op))
		}
		return Int32Value(result), nil

	case TypeInt64:
		i1, err := ValueInt64(v1)
		if err != nil {
			return nil, err
		}
		i2, err := ValueInt64(v2)
		if err != nil {
			return nil, err
		}
		var result int64
		switch b.op {
		case '/':
			result = i1 / i2
		case '*':
			result = i1 * i2
		case '%':
			result = i1 % i2
		case '+':
			result = i1 + i2
		case '-':
			result = i1 - i2
		case TLeftShift:
			result = i1 << i2
		case TRightShift:
			result = i1 >> i2
		default:
			return nil,
				NewError(b.col, fmt.Errorf("unsupport binary operand '%s'",
					b.op))
		}
		return Int64Value(result), nil

	case TypeFloat64:
		i1, err := ValueFloat64(v1)
		if err != nil {
			return nil, err
		}
		i2, err := ValueFloat64(v2)
		if err != nil {
			return nil, err
		}
		var result float64
		switch b.op {
		case '/':
			result = i1 / i2
		case '*':
			result = i1 * i2
		case '+':
			result = i1 + i2
		case '-':
			result = i1 - i2
		default:
			return nil,
				NewError(b.col, fmt.Errorf("unsupport binary operand '%s'",
					b.op))
		}
		return Float64Value(result), nil

	default:
		return nil,
			NewError(b.col,
				fmt.Errorf("unsupport values %s and %s for binary operand '%s'",
					v1, v2, b.op))
	}
}

type unary struct {
	op    TokenType
	col   int
	value Expr
}

func (n unary) String() string {
	return fmt.Sprintf("-%s", n.value)
}

func (n unary) Eval() (Value, error) {
	val, err := n.value.Eval()
	if err != nil {
		return nil, err
	}
	switch val.Type() {
	case TypeInt32:
		ival, err := ValueInt32(val)
		if err != nil {
			return nil, err
		}
		var result int32
		switch n.op {
		case '-':
			result = -ival
		default:
			return nil, NewError(n.col, fmt.Errorf("unsupported %s unary %s",
				val.Type(), n.op))
		}
		return Int32Value(result), nil

	case TypeInt64:
		ival, err := ValueInt64(val)
		if err != nil {
			return nil, err
		}
		var result int64
		switch n.op {
		case '-':
			result = -ival
		default:
			return nil, NewError(n.col, fmt.Errorf("unsupported %s unary %s",
				val.Type(), n.op))
		}
		return Int64Value(result), nil

	case TypeFloat64:
		ival, err := ValueFloat64(val)
		if err != nil {
			return nil, err
		}
		var result float64
		switch n.op {
		case '-':
			result = -ival
		default:
			return nil, NewError(n.col, fmt.Errorf("unsupported %s unary %s",
				val.Type(), n.op))
		}
		return Float64Value(result), nil

	default:
		return nil,
			NewError(n.col, fmt.Errorf("unsupport %s value %s for unary %s",
				val.Type(), val, n.op))
	}
}
