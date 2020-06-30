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

// Input implements command input and output handler.
type Input struct {
	in     *bufio.Reader
	out    io.Writer
	prompt string
	line   []rune
}

// TokenType specifies token types.
type TokenType byte

// Command tokens.
const (
	TIdentifier TokenType = iota
)

// Token specifies command token value.
type Token struct {
	Type   TokenType
	StrVal string
}

func (t *Token) String() string {
	switch t.Type {
	case TIdentifier:
		return t.StrVal

	default:
		return fmt.Sprintf("{Token %d}", t.Type)
	}
}

// NewInput creates a new I/O handler.
func NewInput(in io.Reader, out io.Writer, prompt string) (*Input, error) {
	return &Input{
		in:     bufio.NewReader(in),
		out:    out,
		prompt: prompt,
	}, nil
}

// GetToken returns the next input token.
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
					Type:   TIdentifier,
					StrVal: string(id),
				}, nil
			}
		}
	}

	return &Token{
		Type:   TIdentifier,
		StrVal: string(r),
	}, nil
}

// Rune returns the next input rune.
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

// UngetRune returns the argument rune for input. The next call to
// Rune() will return it instead of consuming input.
func (in *Input) UngetRune(r rune) {
	in.line = append([]rune{r}, in.line...)
}
