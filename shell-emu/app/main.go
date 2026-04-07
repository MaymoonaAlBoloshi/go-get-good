package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"os"
)

func main() {
	commandList := newCommandList()
	var completerItems []readline.PrefixCompleterInterface

	for name := range commandList {
		completerItems = append(completerItems, readline.PcItem(name))
	}

	autoCompleter := readline.NewPrefixCompleter(completerItems...)

	fmt.Print(autoCompleter)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: autoCompleter,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer rl.Close()

	for {
		userInput, err := rl.Readline()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		runCommand(userInput, commandList)
	}
}
