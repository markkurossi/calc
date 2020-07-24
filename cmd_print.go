//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
)

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

func cmdPrint() (int, error) {
	base := Base10

	t, col, err := input.GetToken()
	if err != nil {
		return col, err
	}
	if t.Type == TSlash {
		// Options.
		t, col, err = input.GetToken()
		if err != nil {
			return col, err
		}
		if t.Type != TIdentifier {
			return col, fmt.Errorf("unexpected token '%s'", t)
		}
		switch t.StrVal {
		case "b":
			base = Base2

		case "o":
			base = Base8

		case "x":
			base = Base16

		default:
			return col, fmt.Errorf("unknown option '%s'", t)
		}
	}
	fmt.Printf("base: %s\n", base)

	return 0, fmt.Errorf("print not implemented yet")
}
