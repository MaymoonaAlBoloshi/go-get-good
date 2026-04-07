package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
)

func runCommand(userInput string, commandList map[string]BuiltinFunc) {
	parts, err := shlex.Split(userInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if len(parts) == 0 {
		return
	}

	commandName := parts[0]
	userInput = strings.TrimSpace(userInput)

	output, errOutput, cleanedParts, err := manageOutput(parts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	defer closeIfNotStd(output, os.Stdout)
	defer closeIfNotStd(errOutput, os.Stderr)

	ctx := CommandContext{
		Args:   cleanedParts,
		Stdout: output,
		Stderr: errOutput,
	}

	if cmd, ok := commandList[commandName]; ok {
		cmd(ctx)
		return
	}

	if _, ok := findInPath(commandName); ok {
		args := cleanedParts[1:]

		executable := exec.Command(commandName, args...)
		executable.Stdout = output
		executable.Stderr = errOutput

		executable.Run()

		return
	}

	fmt.Println(userInput + ": command not found")
}

func findInPath(name string) (string, bool) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", false
	}
	return path, true
}

func manageOutput(parts []string) (*os.File, *os.File, []string, error) {
	if len(parts) < 3 {
		return os.Stdout, os.Stderr, parts, nil
	}

	operator := parts[len(parts)-2]
	outputFileName := parts[len(parts)-1]
	cleanedParts := parts[:len(parts)-2]

	isStdErr := operator == "2>" || operator == "2>>"
	isAppend := operator == ">>" || operator == "1>>" || operator == "2>>"
	isRedirect := operator == ">" || operator == "1>" || operator == "2>" || isAppend

	if !isRedirect {
		return os.Stdout, os.Stderr, parts, nil
	}

	var file *os.File
	var err error

	if isAppend {
		file, err = os.OpenFile(outputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(outputFileName)
	}

	if err != nil {
		return nil, nil, nil, err
	}

	if isStdErr {
		return os.Stdout, file, cleanedParts, nil
	}

	return file, os.Stderr, cleanedParts, nil
}

func closeIfNotStd(f *os.File, std *os.File) {
	if f != nil && f != std {
		f.Close()
	}
}
