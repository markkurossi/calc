//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Command defines a built-in command.
type Command struct {
	Name  string
	Title string
	Func  func() error
}

var (
	input    *Input
	commands []Command
)

func init() {
	commands = append(commands, []Command{
		{
			Name:  "help",
			Title: "Print help information",
			Func:  help,
		},
		{
			Name:  "print",
			Title: "Print expression value according to format",
			Func:  cmdPrint,
		},
		{
			Name:  "quit",
			Title: "Exit calc",
			Func: func() error {
				os.Exit(0)
				return nil
			},
		},
	}...)
}

func help() error {
	fmt.Printf("Available commands are:\n\n")
	for _, cmd := range commands {
		fmt.Printf("%s -- %s\n", cmd.Name, cmd.Title)
	}
	return nil
}

func main() {
	flag.Parse()

	log.SetFlags(0)

	var err error

	input, err = NewInput(os.Stdin, os.Stdout, "(calc) ")
	if err != nil {
		log.Fatal(err)
	}

	for {
		t, err := input.GetToken()
		if err != nil {
			log.Fatal(err)
		}
		name := t.String()

		var matches []Command

		for _, cmd := range commands {
			if strings.HasPrefix(cmd.Name, name) {
				matches = append(matches, cmd)
			}
		}
		if len(matches) == 0 {
			log.Printf("Undefined command: \"%s\".  Try \"help\"\n", name)
		} else if len(matches) > 1 {
			fmt.Fprintf(os.Stderr, "Ambiguous command \"%s\":", name)
			for idx, m := range matches {
				if idx > 0 {
					fmt.Fprintf(os.Stderr, ",")
				}
				fmt.Fprintf(os.Stderr, " %s", m.Name)
			}
			fmt.Fprintf(os.Stderr, "\n")
		} else {
			err := matches[0].Func()
			if err != nil {
				col := Column(err)
				if col > 0 {
					var ind string

					for i := 0; i < col; i++ {
						ind += " "
					}
					ind += "^"
					log.Printf("%s\n", ind)
				}
				log.Printf("error: %s\n", err)
			}
			input.FlushEOL()
		}
	}
}
