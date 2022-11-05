//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package main

var Shells = map[string]LazyShell{
	bash: func() Shell {
		return Shell{
			EntryPoint: "bash --login",
			Prompt:     "\\[\\e[38;2;90;86;224m\\]> \\[\\e[0m\\]",
			Command:    ` set +o history; unset PROMPT_COMMAND; export PS1="%s"; clear;`,
		}
	},
	zsh: func() Shell {
		return Shell{
			EntryPoint: "bash --login",
			Prompt:     `%F{#5B56E0}> %F{reset_color}`,
			Command:    ` clear; zsh --login --histnostore; unsetopt PROMPT_SP; unset PROMPT; export PS1="%s"; clear`,
		}
	},
	fish: func() Shell {
		return Shell{
			EntryPoint: "bash --login",
			Prompt:     `function fish_prompt; echo -e "$(set_color 5B56E0)> $(set_color normal)"; end`,
			Command:    `clear; fish --login --private -C 'function fish_greeting; end' -C '%s'`,
		}
	},
	pwsh: func() Shell {
		return Shell{
			EntryPoint: "bash --login",
			Prompt:     "Function prompt {Write-Host \"> \" -ForegroundColor Blue -NoNewLine; return \"`0\" }",
			Command:    ` clear; powershell -NoLogo -NoExit -Command 'Set-PSReadLineOption -HistorySaveStyle SaveNothing; %s'`,
		}
	},
}

var defaultShell = bash
