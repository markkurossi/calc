//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/markkurossi/tabulate"
)

func cmdPrint() error {
	base := Base10
	asCharacter := false

	t, err := input.GetToken()
	if err != nil {
		return err
	}
	if t.Type == TDiv {
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

		case "t":
			base = BaseBinary

		case "c":
			asCharacter = true

		default:
			return NewError(t.Column, fmt.Errorf("unknown option '%s'", t))
		}
	} else {
		input.UngetToken(t)
	}

	expr, err := parseExpr()
	if err != nil {
		return err
	}

	val, err := expr.Eval()
	if err != nil {
		return err
	}

	if asCharacter {
		return printAsCharacter(val)
	}
	fmt.Printf("%s\n", val.Format(Options{
		Base: base,
	}))

	return nil
}

func printAsCharacter(v Value) error {
	r, err := ValueInt32(v)
	if err != nil {
		return err
	}

	tab := tabulate.New(tabulate.Simple)
	tab.Header("Format").SetAlign(tabulate.MR)
	tab.Header("Value").SetAlign(tabulate.ML)

	row := tab.Row()
	row.Column("Decimal")
	row.Column(fmt.Sprintf("%d", r))

	row = tab.Row()
	row.Column("Unicode")
	if r <= 0xffff {
		row.Column(fmt.Sprintf("\\u%04x", r))
	} else {
		row.Column(fmt.Sprintf("\\U%08x", r))
	}

	row = tab.Row()
	row.Column("Symbol")
	if unicode.IsPrint(r) {
		row.Column(fmt.Sprintf("%c", r))
	} else {
		row.Column("{unprintable}")
	}
	tab.Print(os.Stdout)

	return nil
}
