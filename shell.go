package main

// Supported shells of VHS
const (
	bash   = "bash"
	cmdexe = "cmd"
	fish   = "fish"
	pwsh   = "pwsh"
	zsh    = "zsh"
)

// Shell is a type that contains a prompt and the command to set up the shell.
type Shell struct {
	EntryPoint string
	Prompt     string
	Command    string
}

type LazyShell func() Shell
