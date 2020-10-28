//
// Copyright (c) 2019 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"io"
	"testing"
)

type TestReadline struct {
	input []string
}

func (tr *TestReadline) Close() error {
	return nil
}

func (tr *TestReadline) Prompt(prompt string) (string, error) {
	if len(tr.input) == 0 {
		return "", io.EOF
	}
	result := tr.input[0]
	tr.input = tr.input[1:]

	return result, nil
}

func (tr *TestReadline) AppendHistory(item string) {
}

var testReadline = &TestReadline{}

func init() {
	var err error
	input, err = NewInput("(test) ", testReadline)
	if err != nil {
		panic(fmt.Sprintf("failed to init tests: %s", err))
	}
}

type exprTest struct {
	in  string
	out string
}

var exprTests = []exprTest{
	{
		in:  "42",
		out: "42",
	},
	{
		in:  "42.1",
		out: "42.1",
	},
	{
		in:  "42+1",
		out: "43",
	},
	{
		in:  "42.1+1",
		out: "43.1",
	},
	{
		in:  ".1",
		out: "0.1",
	},
	{
		in:  "0.1",
		out: "0.1",
	},
	{
		in:  "1.1",
		out: "1.1",
	},
}

func TestExpr(t *testing.T) {
	for idx, test := range exprTests {
		testReadline.input = []string{test.in}
		expr, err := parseExpr()
		if err != nil {
			t.Errorf("test %d: failed to parse '%s': %s", idx, test.in, err)
			continue
		}
		val, err := expr.Eval()
		if err != nil {
			t.Errorf("test %d: eval failed: %s", idx, err)
			continue
		}
		out := val.String()
		if out != test.out {
			t.Errorf("test %d: unexpected result '%s', expected '%s'",
				idx, out, test.out)
		}
	}
}
