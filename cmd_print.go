//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
)

func cmdPrint() error {
	base := Base10

	t, err := input.GetToken()
	if err != nil {
		return err
	}
	if t.Type == TSlash {
		// Options.
		t, err = input.GetToken()
		if err != nil {
			return err
		}
		if t.Type != TIdentifier {
			return NewError(t.Column, fmt.Errorf("unexpected token '%s'", t))
		}
		switch t.StrVal {
		case "b":
			base = Base2

		case "o":
			base = Base8

		case "x":
			base = Base16

		default:
			return NewError(t.Column, fmt.Errorf("unknown option '%s'", t))
		}
	}

	expr, err := ParseExpr()
	if err != nil {
		return err
	}

	val, err := expr.Eval()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", val.Format(Options{
		Base: base,
	}))

	return fmt.Errorf("print not implemented yet")
}
