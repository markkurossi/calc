//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
)

// Type defines the supported primitive types.
type Type int

// Supported primitive types.
const (
	TypeBool Type = iota
	TypeInt8
	TypeUint8
	TypeInt16
	TypeUint16
	TypeInt32
	TypeUint32
	TypeInt64
	TypeUint64
)

var typeNames = map[Type]string{
	TypeBool:   "bool",
	TypeInt8:   "int8",
	TypeUint8:  "uint8",
	TypeInt16:  "int16",
	TypeUint16: "uint16",
	TypeInt32:  "int32",
	TypeUint32: "uint32",
	TypeInt64:  "int64",
	TypeUint64: "uint64",
}

func (t Type) String() string {
	name, ok := typeNames[t]
	if ok {
		return name
	}
	return fmt.Sprintf("{Type %d}", t)
}

// ConversionType returns the type conversion type for the argument
// values i.e. the smallest type that is capable to represent both
// argument values.
func ConversionType(value1, value2 Value) (Type, error) {
	t1 := value1.Type()
	t2 := value2.Type()
	if t1 > t2 {
		return t1, nil
	}
	return t2, nil
}

// ValueBool returns the value as bool.
func ValueBool(value Value) (bool, error) {
	switch v := value.(type) {
	case BoolValue:
		return bool(v), nil
	case Int8Value:
		if v != 0 {
			return true, nil
		}
		return false, nil
	case Int16Value:
		if v != 0 {
			return true, nil
		}
		return false, nil
	case Int32Value:
		if v != 0 {
			return true, nil
		}
		return false, nil
	case Int64Value:
		if v != 0 {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("type conversion from %T to int failed", value)
}

// ValueInt8 returns the value as int8.
func ValueInt8(value Value) (int8, error) {
	switch v := value.(type) {
	case BoolValue:
		if v {
			return int8(1), nil
		}
		return int8(0), nil
	case Int8Value:
		return int8(v), nil
	case Int16Value:
		return int8(v), nil
	case Int32Value:
		return int8(v), nil
	case Int64Value:
		return int8(v), nil
	}
	return 0, fmt.Errorf("type conversion from %T to int8 failed", value)
}

// ValueInt16 returns the value as int16.
func ValueInt16(value Value) (int16, error) {
	switch v := value.(type) {
	case BoolValue:
		if v {
			return int16(1), nil
		}
		return int16(0), nil
	case Int8Value:
		return int16(v), nil
	case Int16Value:
		return int16(v), nil
	case Int32Value:
		return int16(v), nil
	case Int64Value:
		return int16(v), nil
	}
	return 0, fmt.Errorf("type conversion from %T to int16 failed", value)
}

// ValueInt32 returns the value as int32.
func ValueInt32(value Value) (int32, error) {
	switch v := value.(type) {
	case BoolValue:
		if v {
			return int32(1), nil
		}
		return int32(0), nil
	case Int8Value:
		return int32(v), nil
	case Int16Value:
		return int32(v), nil
	case Int32Value:
		return int32(v), nil
	case Int64Value:
		return int32(v), nil
	}
	return 0, fmt.Errorf("type conversion from %T to int32 failed", value)
}

// ValueInt64 returns the value as int64.
func ValueInt64(value Value) (int64, error) {
	switch v := value.(type) {
	case BoolValue:
		if v {
			return int64(1), nil
		}
		return int64(0), nil
	case Int8Value:
		return int64(v), nil
	case Int16Value:
		return int64(v), nil
	case Int32Value:
		return int64(v), nil
	case Int64Value:
		return int64(v), nil
	}
	return 0, fmt.Errorf("type conversion from %T to int64 failed", value)
}
