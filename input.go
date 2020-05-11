//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Input struct {
	in     *bufio.Reader
	out    io.Writer
	prompt string
	line   []rune
}

type TokenType byte

const (
	T_Identifier TokenType = iota
)

type Token struct {
	Type   TokenType
	StrVal string
}

func (t *Token) String() string {
	switch t.Type {
	case T_Identifier:
		return t.StrVal

	default:
		return fmt.Sprintf("{Token %d}", t.Type)
	}
}

func NewInput(in io.Reader, out io.Writer, prompt string) (*Input, error) {
	return &Input{
		in:     bufio.NewReader(in),
		out:    out,
		prompt: prompt,
	}, nil
}

func (in *Input) GetToken() (*Token, error) {
	var r rune
	var err error

	for {
		r, err = in.Rune()
		if err != nil {
			return nil, err
		}
		if !unicode.IsSpace(r) {
			break
		}
	}

	if unicode.IsLetter(r) {
		id := []rune{r}
		for {
			r, err = in.Rune()
			if err != nil {
				return nil, err
			}
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				id = append(id, r)
			} else {
				in.UngetRune(r)
				return &Token{
					Type:   T_Identifier,
					StrVal: string(id),
				}, nil
			}
		}
	}

	return &Token{
		Type:   T_Identifier,
		StrVal: string(r),
	}, nil
}

func (in *Input) Rune() (rune, error) {
	if len(in.line) == 0 {
		fmt.Fprintf(in.out, "%s", in.prompt)

		line, err := in.in.ReadString('\n')
		if err != nil {
			return 0, err
		}
		in.line = []rune(line)
	}
	r := in.line[0]
	in.line = in.line[1:]
	return r, nil
}

func (in *Input) UngetRune(r rune) {
	in.line = append([]rune{r}, in.line...)
}
