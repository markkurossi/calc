//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"crypto/rand"
	bin "encoding/binary"
	"fmt"
)

var (
	_ Expr = &Builtin{}
)

// Builtin implements builtin functions.
type Builtin struct {
	name string
	col  int
	args []Expr
}

// Eval implements Expr.Eval.
func (bi *Builtin) Eval() (Value, error) {
	switch bi.name {
	case "random":
		if len(bi.args) > 1 {
			return nil, NewError(bi.col,
				fmt.Errorf("%s: too many arguments", bi.name))
		}
		var buf [8]byte
		_, err := rand.Read(buf[:])
		if err != nil {
			return nil, NewError(bi.col, fmt.Errorf("%s: %s", bi.name, err))
		}
		return Int64Value(bin.BigEndian.Uint64(buf[:])), nil

	default:
		return nil, NewError(bi.col, fmt.Errorf("unknown function: '%s'",
			bi.name))
	}
}
