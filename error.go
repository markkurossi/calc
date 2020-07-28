//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

type Error struct {
	Col int
	Err error
}

func (e Error) Error() string {
	return e.Err.Error()
}

func NewError(col int, err error) *Error {
	return &Error{
		Col: col,
		Err: err,
	}
}

func Column(err error) int {
	e, ok := err.(*Error)
	if ok {
		return e.Col
	}
	return 0
}
