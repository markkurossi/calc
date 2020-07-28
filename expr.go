//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
)

type Expr interface {
	Eval() (Value, error)
}

type Value interface {
	String() string
	Format(options Options) string
}

type Options struct {
	Base Base
}

// Base defines the output base for numbers.
type Base int

// Supported output bases.
const (
	Base2 Base = iota
	Base8
	Base10
	Base16
)

var bases = map[Base]string{
	Base2:  "2",
	Base8:  "8",
	Base10: "10",
	Base16: "16",
}

func (b Base) String() string {
	name, ok := bases[b]
	if ok {
		return name
	}
	return fmt.Sprintf("{base %d}", b)
}

func ParseExpr() (Expr, error) {
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
	return parseMultiplicative()
}

func parseMultiplicative() (Expr, error) {
	return parseUnary()
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
	default:
		input.UngetToken(t)
		return nil, NewError(t.Column, fmt.Errorf("unexpected token '%s'", t))
	}
}
