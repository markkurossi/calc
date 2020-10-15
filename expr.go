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
	_ Expr = BoolValue(false)
	_ Expr = Int8Value(0)
	_ Expr = Int16Value(0)
	_ Expr = Int32Value(0)
	_ Expr = Int64Value(0)
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
		case TAdd, TSub:

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
		case TMult, TDiv, TPercent:

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
	return parsePostfix()
}

func parsePostfix() (Expr, error) {
	t, err := input.GetToken()
	if err != nil {
		return nil, err
	}
	switch t.Type {
	case TInteger:
		return t.IntVal, nil

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
		case TDiv:
			result = i1 / i2
		case TMult:
			result = i1 * i2
		case TPercent:
			result = i1 % i2
		case TAdd:
			result = i1 + i2
		case TSub:
			result = i1 - i2
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
		case TDiv:
			result = i1 / i2
		case TMult:
			result = i1 * i2
		case TPercent:
			result = i1 % i2
		case TAdd:
			result = i1 + i2
		case TSub:
			result = i1 - i2
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
		case TDiv:
			result = i1 / i2
		case TMult:
			result = i1 * i2
		case TPercent:
			result = i1 % i2
		case TAdd:
			result = i1 + i2
		case TSub:
			result = i1 - i2
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
		case TDiv:
			result = i1 / i2
		case TMult:
			result = i1 * i2
		case TPercent:
			result = i1 % i2
		case TAdd:
			result = i1 + i2
		case TSub:
			result = i1 - i2
		default:
			return nil,
				NewError(b.col, fmt.Errorf("unsupport binary operand '%s'",
					b.op))
		}
		return Int64Value(result), nil

	default:
		return nil,
			NewError(b.col,
				fmt.Errorf("unsupport values %s and %s for binary operand '%s'",
					v1, v2, b.op))
	}
}
