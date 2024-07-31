//
// Copyright (c) 2020, 2024 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"math/big"
	"strconv"
)

var (
	_ Value = BoolValue(false)
	_ Value = Int8Value(0)
	_ Value = Int16Value(0)
	_ Value = Int32Value(0)
	_ Value = Int64Value(0)
	_ Value = Float64Value(0)
	_ Value = BigFloatValue{
		f: big.NewFloat(0),
	}
)

// Value implements a value.
type Value interface {
	String() string
	Format(options Options) string
	Type() Type
}

// Options define value output options.
type Options struct {
	Base   Base
	String bool
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

var formats = map[Base]byte{
	Base2:      'b',
	Base8:      'f',
	Base10:     'f',
	Base16:     'x',
	BaseBinary: 'b',
}

// FloatFormat returns base format for the strconv.FormatFloat
// function.
func (b Base) FloatFormat() byte {
	fmt, ok := formats[b]
	if ok {
		return fmt
	}
	return 'f'
}

func stringify(v int64, base Base) string {
	var result string

	for i := 56; i >= 0; i -= 8 {
		b := (v >> i) & 0xff
		if b == 0 {
			if len(result) > 0 {
				result += "."
			}
		} else {
			result += string(rune(b))
		}
	}

	return result
}

// BoolValue implements bool values as Value.
type BoolValue bool

func (v BoolValue) String() string {
	if v {
		return "true"
	}
	return "false"
}

// Format implements Value.Format().
func (v BoolValue) Format(options Options) string {
	return v.String()
}

// Type implements Value.Type().
func (v BoolValue) Type() Type {
	return TypeBool
}

// Eval implements Expr.Eval().
func (v BoolValue) Eval() (Value, error) {
	return v, nil
}

// Int8Value implements int values as Value.
type Int8Value int

func (v Int8Value) String() string {
	return strconv.FormatInt(int64(v), 16)
}

// Format implements Value.Format().
func (v Int8Value) Format(options Options) string {
	if options.String {
		return stringify(int64(v), options.Base)
	}
	return options.Base.Prefix() +
		strconv.FormatInt(int64(v), options.Base.Base())
}

// Type implements Value.Type().
func (v Int8Value) Type() Type {
	return TypeInt8
}

// Eval implements Expr.Eval().
func (v Int8Value) Eval() (Value, error) {
	return v, nil
}

// Int16Value implements int values as Value.
type Int16Value int

func (v Int16Value) String() string {
	return strconv.FormatInt(int64(v), 16)
}

// Format implements Value.Format().
func (v Int16Value) Format(options Options) string {
	if options.String {
		return stringify(int64(v), options.Base)
	}
	return options.Base.Prefix() +
		strconv.FormatInt(int64(v), options.Base.Base())
}

// Type implements Value.Type().
func (v Int16Value) Type() Type {
	return TypeInt16
}

// Eval implements Expr.Eval().
func (v Int16Value) Eval() (Value, error) {
	return v, nil
}

// Int32Value implements int values as Value.
type Int32Value int

func (v Int32Value) String() string {
	return strconv.FormatInt(int64(v), 16)
}

// Format implements Value.Format().
func (v Int32Value) Format(options Options) string {
	if options.String {
		return stringify(int64(v), options.Base)
	}
	return options.Base.Prefix() +
		strconv.FormatInt(int64(v), options.Base.Base())
}

// Type implements Value.Type().
func (v Int32Value) Type() Type {
	return TypeInt32
}

// Eval implements Expr.Eval().
func (v Int32Value) Eval() (Value, error) {
	return v, nil
}

// Int64Value implements int64 values as Value.
type Int64Value int64

func (v Int64Value) String() string {
	return strconv.FormatInt(int64(v), 10)
}

// Format implements Value.Format().
func (v Int64Value) Format(options Options) string {
	if options.String {
		return stringify(int64(v), options.Base)
	}
	return options.Base.Prefix() +
		strconv.FormatInt(int64(v), options.Base.Base())
}

// Type implements Value.Type().
func (v Int64Value) Type() Type {
	return TypeInt64
}

// Eval implements Expr.Eval().
func (v Int64Value) Eval() (Value, error) {
	return v, nil
}

// Float64Value implements float64 values as Value.
type Float64Value float64

func (v Float64Value) String() string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

// Format implements Value.Format().
func (v Float64Value) Format(options Options) string {
	if options.String {
		return stringify(int64(v), options.Base)
	}
	return strconv.FormatFloat(float64(v), options.Base.FloatFormat(), -1, 64)
}

// Type implements Value.Type().
func (v Float64Value) Type() Type {
	return TypeFloat64
}

// Eval implements Expr.Eval().
func (v Float64Value) Eval() (Value, error) {
	return v, nil
}

// BigFloatValue implements big.Float values as Value.
type BigFloatValue struct {
	f *big.Float
}

func (v BigFloatValue) String() string {
	return v.f.Text('f', -1)
}

// Format implements Value.Format().
func (v BigFloatValue) Format(options Options) string {
	if options.String {
		ui64, _ := v.f.Uint64()
		return stringify(int64(ui64), options.Base)
	}
	return v.f.Text(options.Base.FloatFormat(), -1)
}

// Type implements Value.Type().
func (v BigFloatValue) Type() Type {
	return TypeBigFloat
}

// Eval implements Expr.Eval().
func (v BigFloatValue) Eval() (Value, error) {
	return v, nil
}
