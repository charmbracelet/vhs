package main

// Supported shells of VHS
const (
	bash       = "bash"
	cmdexe     = "cmd"
	fish       = "fish"
	powershell = "powershell"
	pwsh       = "pwsh"
	zsh        = "zsh"
)

// Shell is a type that contains a prompt and the command to set up the shell.
type Shell struct {
	Prompt  string
	Command string
	Clear   string
}

// Shells contains a mapping from shell names to their Shell struct.
var Shells = map[string]Shell{
	bash: {
		Prompt:  "\\[\\e[38;2;90;86;224m\\]> \\[\\e[0m\\]",
		Command: ` PS1="%s" bash --login --norc --noprofile +o history; clear;`,
		Clear:   "clear",
	},
	zsh: {
		Prompt:  `%F{#5B56E0}> %F{reset_color}`,
		Command: ` PROMPT="%s" zsh --login --histnostore --no-rcs; clear;`,
		Clear:   "clear",
	},
	fish: {
		Prompt:  `function fish_prompt; echo -e "$(set_color 5B56E0)> $(set_color normal)"; end`,
		Command: ` fish --login --no-config --private -C 'function fish_greeting; end' -C '%s'; clear;`,
		Clear:   "clear",
	},
	powershell: {
		Prompt:  `Set-PSReadLineOption -HistorySaveStyle SaveNothing; Function prompt { Write-Host -ForegroundColor Blue -NoNewLine '>'; return ' ' }`,
		Command: ` clear; powershell -Login -NoLogo -NoExit -NoProfile -Command %q`,
		Clear:   "clear",
	},
	pwsh: {
		Prompt:  `Set-PSReadLineOption -HistorySaveStyle SaveNothing; Function prompt { Write-Host -ForegroundColor Blue -NoNewLine '>'; return ' ' }`,
		Command: ` clear; pwsh -Login -NoLogo -NoExit -NoProfile -Command %q`,
		Clear:   "clear",
	},
	cmdexe: {
		Prompt:  "$g",
		Command: ` cls && set prompt=%s && cls`,
		Clear:   "cls",
	},
}
