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
	"strconv"
	"unicode"
)

// Input implements command input and output handler.
type Input struct {
	in     *bufio.Reader
	out    io.Writer
	prompt string
	line   []rune
	col    int
	ungot  *Token
}

// TokenType specifies token types.
type TokenType byte

// Command tokens.
const (
	TIdentifier TokenType = iota
	TInteger
	TDiv
	TMult
	TPercent
	TAdd
	TSub
)

var tokenTypes = map[TokenType]string{
	TIdentifier: "identifier",
	TInteger:    "integer",
	TDiv:        "/",
	TMult:       "*",
	TPercent:    "%",
	TAdd:        "+",
	TSub:        "-",
}

func (t TokenType) String() string {
	name, ok := tokenTypes[t]
	if ok {
		return name
	}
	return fmt.Sprintf("{TokenType %d}", t)
}

// Token specifies command token value.
type Token struct {
	Column int
	Type   TokenType
	StrVal string
	IntVal int64
}

func (t *Token) String() string {
	switch t.Type {
	case TIdentifier:
		return t.StrVal

	case TInteger:
		return fmt.Sprintf("%d", t.IntVal)

	case TDiv, TMult, TPercent, TAdd, TSub:
		return t.Type.String()

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

// FlushEOL discards the current input line.
func (in *Input) FlushEOL() {
	in.ungot = nil
	in.line = []rune{}
}

// HasToken tests if input has any tokens without prompting user.
func (in *Input) HasToken() bool {
	if in.ungot != nil {
		return true
	}

	for in.HasRune() {
		r, _, err := in.Rune()
		if err != nil {
			return false
		}
		if !unicode.IsSpace(r) {
			in.UngetRune(r)
			return true
		}
	}
	return false
}

// GetToken returns the next input token.
func (in *Input) GetToken() (*Token, error) {
	var r rune
	var col, c int
	var err error

	if in.ungot != nil {
		ret := in.ungot
		in.ungot = nil
		return ret, nil
	}

	for {
		r, col, err = in.Rune()
		if err != nil {
			return nil, NewError(col, err)
		}
		if !unicode.IsSpace(r) {
			break
		}
	}
	switch r {
	case '/':
		return &Token{
			Column: col,
			Type:   TDiv,
		}, nil

	case '*':
		return &Token{
			Column: col,
			Type:   TMult,
		}, nil

	case '%':
		return &Token{
			Column: col,
			Type:   TPercent,
		}, nil

	case '+':
		return &Token{
			Column: col,
			Type:   TAdd,
		}, nil

	case '-':
		return &Token{
			Column: col,
			Type:   TSub,
		}, nil

	default:
		if unicode.IsLetter(r) {
			id := []rune{r}
			for {
				r, c, err = in.Rune()
				if err != nil {
					return nil, NewError(c, err)
				}
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
					id = append(id, r)
				} else {
					in.UngetRune(r)
					return &Token{
						Column: col,
						Type:   TIdentifier,
						StrVal: string(id),
					}, nil
				}
			}
		}
		if unicode.IsDigit(r) {
			// XXX 0{x,b,o}DIGITS
			val := []rune{r}
			for {
				r, c, err = in.Rune()
				if err != nil {
					return nil, NewError(c, err)
				}
				if unicode.IsDigit(r) {
					val = append(val, r)
				} else {
					in.UngetRune(r)
					i64, err := strconv.ParseInt(string(val), 10, 64)
					if err != nil {
						return nil, NewError(col, err)
					}
					return &Token{
						Column: col,
						Type:   TInteger,
						IntVal: i64,
					}, nil
				}
			}
		}
		return nil, NewError(col, fmt.Errorf("unexpected character '%c'", r))
	}
}

// UngetToken ungets the token. The next call to GetToken will returns
// the token instead of consuming input stream.
func (in *Input) UngetToken(t *Token) {
	in.ungot = t
}

// HasRune tests if input is not empty.
func (in *Input) HasRune() bool {
	return len(in.line) > 0
}

// Rune returns the next input rune.
func (in *Input) Rune() (rune, int, error) {
	if len(in.line) == 0 {
		fmt.Fprintf(in.out, "%s", in.prompt)

		line, err := in.in.ReadString('\n')
		if err != nil {
			return 0, len(in.prompt), err
		}
		in.line = []rune(line)
		in.col = len(in.prompt)
	}
	r := in.line[0]
	in.line = in.line[1:]
	in.col++
	return r, in.col - 1, nil
}

// UngetRune returns the argument rune for input. The next call to
// Rune() will return it instead of consuming input.
func (in *Input) UngetRune(r rune) {
	in.line = append([]rune{r}, in.line...)
	in.col--
}
