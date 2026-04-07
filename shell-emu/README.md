# shell-emu

Tiny shell emulator in Go.

This was written for learning purposes, not as a serious shell. The code is pretty bad in places on purpose. It is messy, incomplete, and a bit naive, because the point was to learn by building it.

## What exists in it

- a simple `$` prompt
- built-in commands: `echo`, `pwd`, `cd`, `type`, `exit`
- external command execution from `PATH`
- basic autocomplete for built-ins
- simple output redirection

## Run it

```bash
go run ./app
```

## Commands

`echo`
Prints text.

```bash
echo hello world
```

`pwd`
Prints the current directory.

```bash
pwd
```

`cd`
Changes directory. `cd` and `cd ~` go home.

```bash
cd /tmp
cd
cd ~
```

`type`
Shows whether something is a built-in command or found in `PATH`.

```bash
type echo
type go
```

`exit`
Closes the shell.

```bash
exit
```

## Redirection

Supported:

```bash
echo hi > out.txt
echo hi >> out.txt
some-command 2> err.txt
some-command 2>> err.txt
```

## Limits

- not a real shell
- no pipes
- no `&&` or `||`
- no env var expansion
- autocomplete only works for built-ins
- redirection only works at the end of the command
