package main

func getPrompt(vhs *VHS) Prompt {
	switch vhs.Options.Shell {
	case "", "bash":
		return bash{}
	case "zsh":
		return zsh{}
	case "fish":
		return fish{}
	case "pwsh":
		return pwsh{}
	default:
		return generic{}
	}
}

// Prompt defines how prompts set themselves up.
type Prompt interface {
	Setup(vhs *VHS)
}

var (
	_ Prompt = bash{}
	_ Prompt = zsh{}
	_ Prompt = fish{}
	_ Prompt = pwsh{}
	_ Prompt = generic{}
)

type (
	bash    struct{}
	zsh     struct{}
	fish    struct{}
	pwsh    struct{}
	generic struct{}
)

func (bash) Setup(vhs *VHS) {
	prompt := vhs.Options.Prompt
	if prompt == "" {
		prompt = "\\[\\e[38;2;90;86;224m\\]> \\[\\e[0m\\]"
	}
	vhs.runShellCommandf(` set +o history; unset PROMPT_COMMAND; export PS1="%s"; clear;`, prompt)
}

func (zsh) Setup(vhs *VHS) {
	prompt := vhs.Options.Prompt
	if prompt == "" {
		prompt = `%F{blue bright dim}> %F{reset_color}`
	}
	vhs.runShellCommandf(" clear; zsh --login --histnostore")
	// PROMPT_SP: read about PROMPT_EOL_MARK
	vhs.runShellCommandf(` unsetopt PROMPT_SP; export PS1="%s"; clear`, prompt)
}

func (fish) Setup(vhs *VHS) {
	prompt := vhs.Options.Prompt
	if prompt == "" {
		prompt = `function fish_prompt; echo -e "$(set_color --dim brblue)> $(set_color normal)"; end`
	}
	noGreeting := "function fish_greeting; end"
	vhs.runShellCommandf(` clear; fish --login --private -C '%s' -C '%s'`, noGreeting, prompt)
}

func (pwsh) Setup(vhs *VHS) {
	prompt := vhs.Options.Prompt
	if prompt == "" {
		prompt = "Function prompt {Write-Host \"> \" -ForegroundColor Blue -NoNewLine; return \"`0\" }"
	}
	vhs.runShellCommandf(` clear; pwsh -Login -NoLogo -NoExit -Command 'Set-PSReadLineOption -HistorySaveStyle SaveNothing; %s'`, prompt)
	// XXX: here would be a great place to reuse whatever we do for the Exec thing, so we can wait for the shell to load, maybe...
	vhs.runShellCommandf(`clear; sleep 1; clear`)
}

func (generic) Setup(vhs *VHS) {
	// XXX: what should we do with prompt here?
	vhs.runShellCommandf(` clear; %s`, vhs.Options.Shell)
}
