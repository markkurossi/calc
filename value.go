//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"strconv"
)

var (
	_ Value = IntegerValue(0)
)

// Value implements a value.
type Value interface {
	String() string
	Format(options Options) string
}

// Options define value output options.
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
	BaseBinary
)

var baseNames = map[Base]string{
	Base2:      "2",
	Base8:      "8",
	Base10:     "10",
	Base16:     "16",
	BaseBinary: "2",
}

func (b Base) String() string {
	name, ok := baseNames[b]
	if ok {
		return name
	}
	return fmt.Sprintf("{base %d}", b)
}

var basePrefixes = map[Base]string{
	Base2:      "0b",
	Base8:      "0",
	Base10:     "",
	Base16:     "0x",
	BaseBinary: "",
}

// Prefix returns the base prefix string.
func (b Base) Prefix() string {
	prefix, ok := basePrefixes[b]
	if ok {
		return prefix
	}
	return ""
}

var bases = map[Base]int{
	Base2:      2,
	Base8:      8,
	Base10:     10,
	Base16:     16,
	BaseBinary: 2,
}

// Base returns the base as integer number.
func (b Base) Base() int {
	base, ok := bases[b]
	if ok {
		return base
	}
	return 10
}

// IntegerValue implements int64 values as Value.
type IntegerValue int64

func (v IntegerValue) String() string {
	return strconv.FormatInt(int64(v), 10)
}

// Format implements Value.Format().
func (v IntegerValue) Format(options Options) string {
	return options.Base.Prefix() +
		strconv.FormatInt(int64(v), options.Base.Base())
}

// Eval implements Expr.Eval().
func (v IntegerValue) Eval() (Value, error) {
	return v, nil
}
