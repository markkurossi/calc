//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
)

// ConversionType returns the type conversion type for the argument
// values
func ConversionType(value1, value2 Value) (Type, error) {
	switch value1.(type) {
	case IntValue:
		switch value2.(type) {
		case IntValue, Int8Value, Int16Value, Int32Value:
			return TypeInt, nil

		case Int64Value:
			return TypeInt64, nil
		}

	case Int8Value:
		switch value2.(type) {
		case IntValue:
			return TypeInt, nil

		case Int8Value:
			return TypeInt8, nil

		case Int16Value:
			return TypeInt16, nil

		case Int32Value:
			return TypeInt32, nil

		case Int64Value:
			return TypeInt64, nil
		}

	case Int16Value:
		switch value2.(type) {
		case IntValue:
			return TypeInt, nil

		case Int8Value, Int16Value:
			return TypeInt16, nil

		case Int32Value:
			return TypeInt32, nil

		case Int64Value:
			return TypeInt64, nil
		}

	case Int32Value:
		switch value2.(type) {
		case IntValue:
			return TypeInt, nil

		case Int8Value, Int16Value, Int32Value:
			return TypeInt32, nil

		case Int64Value:
			return TypeInt64, nil
		}

	case Int64Value:
		switch value2.(type) {
		case IntValue, Int8Value, Int16Value, Int32Value, Int64Value:
			return TypeInt64, nil
		}
	}

	return 0, fmt.Errorf("type conversion failed for types %T and %T",
		value1, value2)
}

// ValueInt returns the value as int.
func ValueInt(value Value) (int, error) {
	switch v := value.(type) {
	case IntValue:
		return int(v), nil
	case Int8Value:
		return int(v), nil
	case Int16Value:
		return int(v), nil
	case Int32Value:
		return int(v), nil
	case Int64Value:
		return int(v), nil
	}
	return 0, fmt.Errorf("type conversion from %T to int failed", value)
}

// ValueInt8 returns the value as int8.
func ValueInt8(value Value) (int8, error) {
	switch v := value.(type) {
	case IntValue:
		return int8(v), nil
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
	case IntValue:
		return int16(v), nil
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
	case IntValue:
		return int32(v), nil
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
	case IntValue:
		return int64(v), nil
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
