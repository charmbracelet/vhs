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
	Command []string
	Env     []string
}

// Shells contains a mapping from shell names to their Shell struct.
var Shells = map[string]Shell{
	bash: {
		Env:     []string{"PS1=\\[\\e[38;2;90;86;224m\\]> \\[\\e[0m\\]", "BASH_SILENCE_DEPRECATION_WARNING=1"},
		Command: []string{"bash", "--noprofile", "--norc", "--login", "+o", "history"},
	},
	zsh: {
		Env:     []string{`PROMPT=%F{#5B56E0}> %F{reset_color}`},
		Command: []string{"zsh", "--histnostore", "--no-rcs"},
	},
	fish: {
		Command: []string{
			"fish",
			"--login",
			"--no-config",
			"--private",
			"-C", "function fish_greeting; end",
			"-C", `function fish_prompt; echo -e "$(set_color 5B56E0)> $(set_color normal)"; end`,
		},
	},
	powershell: {
		Command: []string{
			"powershell",
			"-Login",
			"-NoLogo",
			"-NoExit",
			"-NoProfile",
			"-Command",
			`Set-PSReadLineOption -HistorySaveStyle SaveNothing; Function prompt { Write-Host -ForegroundColor Blue -NoNewLine '>'; return ' ' }`,
		},
	},
	pwsh: {
		Command: []string{
			"pwsh",
			"-Login",
			"-NoLogo",
			"-NoExit",
			"-NoProfile",
			"-Command",
			`Set-PSReadLineOption -HistorySaveStyle SaveNothing; Function prompt { Write-Host -ForegroundColor Blue -NoNewLine '>'; return ' ' }`,
		},
	},
	cmdexe: {
		Command: []string{"cmd.exe", "/k", "prompt=^> "},
	},
}
