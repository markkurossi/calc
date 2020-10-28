//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

// Readline implements line-based user input.
type Readline interface {
	Close() error
	Prompt(prompt string) (string, error)
	AppendHistory(item string)
}
