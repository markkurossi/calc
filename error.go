//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

// Error implements error with input location information.
type Error struct {
	Col int
	Err error
}

func (e Error) Error() string {
	return e.Err.Error()
}

// NewError creates a new error from the error and location
// information.
func NewError(col int, err error) *Error {
	return &Error{
		Col: col,
		Err: err,
	}
}

// Column returns the error location information.
func Column(err error) int {
	e, ok := err.(*Error)
	if ok {
		return e.Col
	}
	return 0
}
