//
// Copyright (c) 2020-2024 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"math/big"
	"strconv"
	"unicode"
)

// Input implements command input and output handler.
type Input struct {
	prompt   string
	line     []rune
	col      int
	ungot    *Token
	readline Readline
}

// TokenType specifies token types.
type TokenType int

// Command tokens.
const (
	TIdentifier TokenType = iota + 256
	TInteger
	TFloat
	TLeftShift
	TRightShift
)

var tokenTypes = map[TokenType]string{
	TIdentifier: "identifier",
	TInteger:    "integer",
	TFloat:      "float",
	TLeftShift:  "<<",
	TRightShift: ">>",
}

func (t TokenType) String() string {
	name, ok := tokenTypes[t]
	if ok {
		return name
	}
	if t < TIdentifier {
		return fmt.Sprintf("%c", rune(t))
	}
	return fmt.Sprintf("{TokenType %d}", t)
}

// Token specifies command token value.
type Token struct {
	Column   int
	Type     TokenType
	StrVal   string
	IntVal   Expr
	FloatVal Expr
}

func (t *Token) String() string {
	switch t.Type {
	case TIdentifier:
		return t.StrVal

	case TInteger:
		return fmt.Sprintf("%v", t.IntVal)

	case TFloat:
		return fmt.Sprintf("%v", t.FloatVal)

	default:
		return t.Type.String()
	}
}

// NewInput creates a new I/O handler.
func NewInput(prompt string, readline Readline) (*Input, error) {
	return &Input{
		prompt:   prompt,
		readline: readline,
	}, nil
}

// Close closes the input.
func (in *Input) Close() {
	in.readline.Close()
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
		r, _, err := in.Rune(false)
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

// GetFirstToken returns the next input token which is the first token
// of a command.
func (in *Input) GetFirstToken() (*Token, error) {
	return in.getToken(true)
}

// GetToken returns the next input token which is a subsequent token
// in a command.
func (in *Input) GetToken() (*Token, error) {
	return in.getToken(false)
}

func (in *Input) getToken(first bool) (*Token, error) {
	var r rune
	var col, c int
	var err error

	if in.ungot != nil {
		ret := in.ungot
		in.ungot = nil
		return ret, nil
	}

	for {
		r, col, err = in.Rune(first)
		if err != nil {
			return nil, NewError(col, err)
		}
		if !unicode.IsSpace(r) {
			break
		}
	}
	switch r {
	case '/', '*', '%', '+', '-', '(', ')', ',':
		return &Token{
			Column: col,
			Type:   TokenType(r),
		}, nil

	case '<':
		n, _, err := in.Rune(first)
		if err != nil {
			return nil, NewError(col, err)
		}
		switch n {
		case '<':
			return &Token{
				Column: col,
				Type:   TLeftShift,
			}, nil

		default:
			in.UngetRune(n)
			return &Token{
				Column: col,
				Type:   TokenType(r),
			}, nil
		}

	case '\'':
		ch, chCol, err := in.Rune(first)
		if err != nil {
			return nil, NewError(col, err)
		}
		if ch == '\\' {
			ch, chCol, err = in.Rune(first)
			if err != nil {
				return nil, NewError(col, err)
			}
			switch ch {
			case 'a':
				ch = '\a'
			case 'b':
				ch = '\b'
			case 'f':
				ch = '\f'
			case 'n':
				ch = '\n'
			case 'r':
				ch = '\r'
			case 't':
				ch = '\t'
			case 'v':
				ch = '\v'
			case '\\':
				ch = '\\'
			case '\'':
				ch = '\''
			default:
				return nil, NewError(chCol,
					fmt.Errorf("unexpected character '%c' in char literal", ch))
			}
		}
		r, col, err = in.Rune(first)
		if err != nil {
			return nil, NewError(col, err)
		}
		if r != '\'' {
			return nil, NewError(chCol,
				fmt.Errorf("unexpected character '%c' in char literal", r))
		}
		return &Token{
			Column: chCol,
			Type:   TInteger,
			IntVal: Int8Value(ch),
		}, nil

	case '.':
		r, c, err = in.Rune(first)
		if err != nil {
			return nil, NewError(c, err)
		}
		if unicode.IsDigit(r) {
			val := []rune{'.', r}
			for {
				r, c, err := in.Rune(first)
				if err != nil {
					return nil, NewError(c, err)
				}
				if unicode.IsDigit(r) {
					val = append(val, r)
				} else {
					in.UngetRune(r)
					break
				}
			}
			return in.parseFloatLiteral(col, val, '.')
		}
		in.UngetRune(r)
		return &Token{
			Column: col,
			Type:   TokenType('.'),
		}, nil

	case '0':
		r, c, err = in.Rune(first)
		if err != nil {
			return nil, NewError(c, err)
		}
		var i64 int64
		switch r {
		case 'b', 'B':
			i64, err = in.readBinaryLiteral([]rune{'0', r})
		case 'o', 'O':
			i64, err = in.readOctalLiteral([]rune{'0', r})
		case 'x', 'X':
			i64, err = in.readHexLiteral([]rune{'0', r})
		case '0', '1', '2', '3', '4', '5', '6', '7':
			i64, err = in.readOctalLiteral([]rune{'0', r})
		case '.':
			val := []rune{'0', r}
			for {
				r, c, err := in.Rune(first)
				if err != nil {
					return nil, NewError(c, err)
				}
				if unicode.IsDigit(r) {
					val = append(val, r)
				} else {
					in.UngetRune(r)
					break
				}
			}
			return in.parseFloatLiteral(col, val, '.')
		default:
			in.UngetRune(r)
		}
		if err != nil {
			return nil, NewError(col, err)
		}
		return &Token{
			Column: col,
			Type:   TInteger,
			IntVal: Int64Value(i64),
		}, nil

	default:
		if unicode.IsLetter(r) {
			id := []rune{r}
			for {
				r, c, err = in.Rune(first)
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
			val := []rune{r}
			var numComma, lastComma, numPeriod, lastPeriod int
			for {
				r, c, err = in.Rune(first)
				if err != nil {
					return nil, NewError(c, err)
				}
				if r == ' ' {
				} else if r == '.' {
					numPeriod++
					lastPeriod = len(val)
					val = append(val, r)
				} else if r == ',' {
					numComma++
					lastComma = len(val)
					val = append(val, r)
				} else if unicode.IsDigit(r) {
					val = append(val, r)
				} else {
					in.UngetRune(r)
					break
				}
			}

			if numPeriod == 1 {
				if numComma == 0 || lastPeriod > lastComma {
					return in.parseFloatLiteral(col, val, '.')
				} else if numComma == 1 && lastComma > lastPeriod {
					return in.parseFloatLiteral(col, val, ',')
				}
			} else if numComma == 1 {
				if numPeriod == 0 || lastComma > lastPeriod {
					return in.parseFloatLiteral(col, val, ',')
				} else if numPeriod == 1 && lastPeriod > lastComma {
					return in.parseFloatLiteral(col, val, '.')
				}
			}
			if numComma > 0 || numPeriod > 0 {
				return nil, NewError(col,
					fmt.Errorf("invalid float number: %v", string(val)))
			}
			i64, err := strconv.ParseInt(string(val), 10, 64)
			if err != nil {
				return nil, NewError(col, err)
			}
			return &Token{
				Column: col,
				Type:   TInteger,
				IntVal: Int64Value(i64),
			}, nil
		}
		return nil, NewError(col, fmt.Errorf("unexpected character '%c'", r))
	}
}

func (in *Input) readBinaryLiteral(val []rune) (int64, error) {
	for {
		r, c, err := in.Rune(false)
		if err != nil {
			return 0, NewError(c, err)
		}
		switch r {
		case '0', '1':
			val = append(val, r)
		default:
			in.UngetRune(r)
			return strconv.ParseInt(string(val), 0, 64)
		}
	}
}

func (in *Input) readOctalLiteral(val []rune) (int64, error) {
	for {
		r, c, err := in.Rune(false)
		if err != nil {
			return 0, NewError(c, err)
		}
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7':
			val = append(val, r)
		default:
			in.UngetRune(r)
			return strconv.ParseInt(string(val), 0, 64)
		}
	}
}

func (in *Input) readHexLiteral(val []rune) (int64, error) {
	for {
		r, c, err := in.Rune(false)
		if err != nil {
			return 0, NewError(c, err)
		}
		if unicode.Is(unicode.Hex_Digit, r) {
			val = append(val, r)
		} else {
			in.UngetRune(r)
			return strconv.ParseInt(string(val), 0, 64)
		}
	}
}

func (in *Input) parseFloatLiteral(col int, val []rune, sep rune) (
	*Token, error) {

	var clean []rune
	for _, r := range val {
		if r == sep {
			clean = append(clean, '.')
		} else if unicode.IsDigit(r) {
			clean = append(clean, r)
		}
	}
	f, _, err := big.ParseFloat(string(clean), 10, 1024, big.ToNearestEven)
	if err != nil {
		return nil, err
	}
	return &Token{
		Column: col,
		Type:   TFloat,
		FloatVal: BigFloatValue{
			f: f,
		},
	}, nil
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
func (in *Input) Rune(first bool) (rune, int, error) {
	if len(in.line) == 0 {
		var prompt = in.prompt
		if !first {
			prompt = "> "
		}

		line, err := in.readline.Prompt(prompt)
		if err != nil {
			return 0, len(in.prompt), err
		}
		in.readline.AppendHistory(line)
		in.line = append([]rune(line), '\n')
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
