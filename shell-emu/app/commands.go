package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type CommandContext struct {
	Args   []string
	Stdout io.Writer
	Stderr io.Writer
}

type BuiltinFunc func(CommandContext)

func newCommandList() map[string]BuiltinFunc {
	var commandList map[string]BuiltinFunc

	commandList = map[string]BuiltinFunc{
		"type": func(ctx CommandContext) {
			if len(ctx.Args) < 2 {
				fmt.Fprintln(ctx.Stdout, "type: missing argument")
				return
			}

			program := ctx.Args[1]
			_, exists := commandList[program]

			if exists {
				fmt.Fprintln(ctx.Stdout, program+" is a shell builtin")
			} else if path, ok := findInPath(program); ok {
				fmt.Fprintln(ctx.Stdout, program+" is "+path)
			} else {
				fmt.Fprintln(ctx.Stdout, program+": not found")
			}
		},
		"echo": func(ctx CommandContext) {
			if len(ctx.Args) == 1 {
				fmt.Fprintln(ctx.Stdout)
				return
			}

			fmt.Fprintln(ctx.Stdout, strings.Join(ctx.Args[1:], " "))
		},
		"pwd": func(ctx CommandContext) {
			pwd, _ := os.Getwd()
			fmt.Fprintln(ctx.Stdout, pwd)
		},
		"cd": func(ctx CommandContext) {
			home, _ := os.UserHomeDir()
			if len(ctx.Args) == 1 || ctx.Args[1] == "~" {
				os.Chdir(home)
				return
			}

			path := ctx.Args[1]
			if err := os.Chdir(path); err != nil {
				fmt.Fprintln(ctx.Stdout, "cd: "+path+": No such file or directory")
				return
			}
		},
		"exit": func(ctx CommandContext) {
			os.Exit(0)
		},
	}

	return commandList
}
