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
	_ Value = Int64Value(0)
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

// Type defines the supported primitive types.
type Type int

// Supported primitive types.
const (
	TypeInt Type = iota
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt64
	TypeUint
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint64
)

var typeNames = map[Type]string{
	TypeInt:    "int",
	TypeInt8:   "int8",
	TypeInt16:  "int16",
	TypeInt32:  "int32",
	TypeInt64:  "int64",
	TypeUint:   "uint",
	TypeUint8:  "uint8",
	TypeUint16: "uint16",
	TypeUint32: "uint32",
	TypeUint64: "uint64",
}

func (t Type) String() string {
	name, ok := typeNames[t]
	if ok {
		return name
	}
	return fmt.Sprintf("{Type %d}", t)
}

// IntValue implements int values as Value.
type IntValue int

func (v IntValue) String() string {
	return strconv.FormatInt(int64(v), 16)
}

// Format implements Value.Format().
func (v IntValue) Format(options Options) string {
	return options.Base.Prefix() +
		strconv.FormatInt(int64(v), options.Base.Base())
}

// Eval implements Expr.Eval().
func (v IntValue) Eval() (Value, error) {
	return v, nil
}

// Int8Value implements int values as Value.
type Int8Value int

func (v Int8Value) String() string {
	return strconv.FormatInt(int64(v), 16)
}

// Format implements Value.Format().
func (v Int8Value) Format(options Options) string {
	return options.Base.Prefix() +
		strconv.FormatInt(int64(v), options.Base.Base())
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
	return options.Base.Prefix() +
		strconv.FormatInt(int64(v), options.Base.Base())
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
	return options.Base.Prefix() +
		strconv.FormatInt(int64(v), options.Base.Base())
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
	return options.Base.Prefix() +
		strconv.FormatInt(int64(v), options.Base.Base())
}

// Eval implements Expr.Eval().
func (v Int64Value) Eval() (Value, error) {
	return v, nil
}
