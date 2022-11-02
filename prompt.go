package main

// Supported shells of VHS
const (
	BASH       = "bash"
	CMD        = "cmd"
	FISH       = "fish"
	POWERSHELL = "pwsh"
	ZSH        = "zsh"
)

// Shell is a type that contains a prompt and the command to set up the shell.
type Shell struct {
	Prompt  string
	Command string
}

// Shells contains a mapping from shell names to their Shell struct.
var Shells = map[string]Shell{
	BASH: {
		Prompt:  "\\[\\e[38;2;90;86;224m\\]> \\[\\e[0m\\]",
		Command: ` set +o history; unset PROMPT_COMMAND; export PS1="%s"; clear;`,
	},
	ZSH: {
		Prompt:  `%F{#5B56E0}> %F{reset_color}`,
		Command: ` clear; zsh --login --histnostore; unsetopt PROMPT_SP; unset PROMPT; export PS1="%s"; clear`,
	},
	FISH: {
		Prompt:  `function fish_prompt; echo -e "$(set_color 5B56E0)> $(set_color normal)"; end`,
		Command: `clear; fish --login --private -C 'function fish_greeting; end' -C '%s'`,
	},
	POWERSHELL: {
		Prompt:  "Function prompt {Write-Host \"> \" -ForegroundColor Blue -NoNewLine; return \"`0\" }",
		Command: ` clear; pwsh -Login -NoLogo -NoExit -Command 'Set-PSReadLineOption -HistorySaveStyle SaveNothing; %s'`,
	},
	CMD: {
		Prompt:  "> ",
		Command: ` clear; export PS1="%s"; clear`,
	},
}
