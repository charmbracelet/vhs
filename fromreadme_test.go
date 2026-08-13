package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFindREADME(t *testing.T) {
	t.Run("finds README.md", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "README.md")
		requireNoErr(t, os.WriteFile(path, []byte("# Hello"), 0o600))

		got, err := findREADME(dir)
		requireNoErr(t, err)
		if got != path {
			t.Errorf("want %q, got %q", path, got)
		}
	})

	t.Run("finds README without extension", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "README")
		requireNoErr(t, os.WriteFile(path, []byte("# Hello"), 0o600))

		got, err := findREADME(dir)
		requireNoErr(t, err)
		if got != path {
			t.Errorf("want %q, got %q", path, got)
		}
	})

	t.Run("prefers README.md over README", func(t *testing.T) {
		dir := t.TempDir()
		requireNoErr(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("# A"), 0o600))
		requireNoErr(t, os.WriteFile(filepath.Join(dir, "README"), []byte("# B"), 0o600))

		got, err := findREADME(dir)
		requireNoErr(t, err)
		want := filepath.Join(dir, "README.md")
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})

	t.Run("error when no README found", func(t *testing.T) {
		dir := t.TempDir()
		_, err := findREADME(dir)
		requireErr(t, err)
	})
}

func TestExtractShellCommands(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name: "extracts from bash block",
			content: "```bash\necho hello\n```",
			want: []string{"echo hello"},
		},
		{
			name: "extracts from sh block",
			content: "```sh\nls -la\n```",
			want: []string{"ls -la"},
		},
		{
			name: "extracts from zsh block",
			content: "```zsh\npwd\n```",
			want: []string{"pwd"},
		},
		{
			name: "extracts from shell block",
			content: "```shell\ndate\n```",
			want: []string{"date"},
		},
		{
			name: "skips non-shell code blocks",
			content: "```go\nfmt.Println(\"hi\")\n```\n```bash\necho hi\n```",
			want: []string{"echo hi"},
		},
		{
			name: "skips comment lines inside code blocks",
			content: "```bash\n# This is a comment\necho hello\n```",
			want: []string{"echo hello"},
		},
		{
			name: "skips empty lines inside code blocks",
			content: "```bash\n\necho hello\n\n```",
			want: []string{"echo hello"},
		},
		{
			name:    "returns nil when no shell blocks",
			content: "```python\nprint('hi')\n```",
			want:    nil,
		},
		{
			name: "extracts from multiple blocks",
			content: "```bash\necho first\n```\n\nSome text.\n\n```sh\necho second\n```",
			want: []string{"echo first", "echo second"},
		},
		{
			name: "extracts multiple lines from one block",
			content: "```bash\necho one\necho two\necho three\n```",
			want: []string{"echo one", "echo two", "echo three"},
		},
		{
			name: "handles tilde fences",
			content: "~~~bash\necho hello\n~~~",
			want: []string{"echo hello"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := extractShellCommands(tc.content, "")
			requireNoErr(t, err)
			if !stringSliceEqual(got, tc.want) {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestExtractShellCommandsWithSection(t *testing.T) {
	tests := []struct {
		name    string
		content string
		section string
		want    []string
	}{
		{
			name: "extracts only from matching section",
			content: `## Usage

` + "```bash\necho usage\n```" + `

## Other

` + "```bash\necho other\n```",
			section: "Usage",
			want:    []string{"echo usage"},
		},
		{
			name: "section match is case-insensitive",
			content: `## INSTALLATION

` + "```bash\necho install\n```",
			section: "installation",
			want:    []string{"echo install"},
		},
		{
			name: "section ends at same-level heading",
			content: `## Usage

` + "```bash\necho usage\n```" + `

## Install

` + "```bash\necho install\n```",
			section: "Usage",
			want:    []string{"echo usage"},
		},
		{
			name: "section ends at higher-level heading",
			content: `### Usage

` + "```bash\necho usage\n```" + `

## Top Level

` + "```bash\necho top\n```",
			section: "Usage",
			want:    []string{"echo usage"},
		},
		{
			name: "subsection included within parent section",
			content: `## Usage

` + "```bash\necho usage\n```" + `

### Advanced

` + "```bash\necho advanced\n```" + `

## Other

` + "```bash\necho other\n```",
			section: "Usage",
			want:    []string{"echo usage", "echo advanced"},
		},
		{
			name: "no match returns nil",
			content: `## Usage

` + "```bash\necho usage\n```",
			section: "Install",
			want:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := extractShellCommands(tc.content, tc.section)
			requireNoErr(t, err)
			if !stringSliceEqual(got, tc.want) {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestFilterByCommand(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		command  string
		want     []string
	}{
		{
			name:     "exact first token match",
			commands: []string{"myapp run", "other run"},
			command:  "myapp",
			want:     []string{"myapp run"},
		},
		{
			name:     "match after sudo prefix",
			commands: []string{"sudo myapp install", "echo hi"},
			command:  "myapp",
			want:     []string{"sudo myapp install"},
		},
		{
			name:     "match after npx prefix",
			commands: []string{"npx myapp build"},
			command:  "myapp",
			want:     []string{"npx myapp build"},
		},
		{
			name:     "match after env prefix",
			commands: []string{"env myapp run"},
			command:  "myapp",
			want:     []string{"env myapp run"},
		},
		{
			name:     "no match returns nil",
			commands: []string{"echo hello", "ls -la"},
			command:  "myapp",
			want:     nil,
		},
		{
			name:     "empty command returns all lines",
			commands: []string{"echo hello", "ls -la"},
			command:  "",
			want:     []string{"echo hello", "ls -la"},
		},
		{
			name:     "basename matching for path-prefixed commands",
			commands: []string{"/usr/local/bin/myapp run"},
			command:  "myapp",
			want:     []string{"/usr/local/bin/myapp run"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := filterByCommand(tc.commands, tc.command)
			if !stringSliceEqual(got, tc.want) {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestTapeDuration(t *testing.T) {
	tests := []struct {
		input time.Duration
		want  string
	}{
		{0, "0s"},
		{500 * time.Millisecond, "500ms"},
		{75 * time.Millisecond, "75ms"},
		{1 * time.Second, "1s"},
		{30 * time.Second, "30s"},
		{1 * time.Minute, "1m"},
		{2 * time.Minute, "2m"},
		{90 * time.Second, "90s"},      // 1m30s — not a whole minute, falls to seconds
		{1500 * time.Millisecond, "1500ms"}, // not a whole second
	}
	for _, tc := range tests {
		got := tapeDuration(tc.input)
		if got != tc.want {
			t.Errorf("tapeDuration(%v): want %q, got %q", tc.input, tc.want, got)
		}
	}
}

func TestBuildFromReadmeTape(t *testing.T) {
	opts := fromReadmeOptions{
		output:      "out.gif",
		fontSize:    15,
		width:       1600,
		height:      900,
		typingSpeed: 75 * time.Millisecond,
		pause:       2 * time.Second,
		waitTimeout: 2 * time.Minute,
		waitPattern: defaultFromReadmeWaitPattern,
	}

	t.Run("header contains settings", func(t *testing.T) {
		tape := buildFromReadmeTape([]string{"echo hi"}, opts)
		for _, want := range []string{
			"Output out.gif",
			"Set FontSize 15",
			"Set Width 1600",
			"Set Height 900",
			"Set TypingSpeed 75ms",
			"Set WaitTimeout 2m",
			"Set WaitPattern",
		} {
			if !strings.Contains(tape, want) {
				t.Errorf("tape missing %q\ntape:\n%s", want, tape)
			}
		}
	})

	t.Run("command produces Type/Sleep/Enter/Wait/Sleep sequence", func(t *testing.T) {
		tape := buildFromReadmeTape([]string{"echo hello"}, opts)
		for _, want := range []string{
			`Type "echo hello"`,
			"Sleep 500ms",
			"Enter",
			"Wait",
			"Sleep 2s",
		} {
			if !strings.Contains(tape, want) {
				t.Errorf("tape missing %q\ntape:\n%s", want, tape)
			}
		}
	})

	t.Run("command with double quotes uses single-quote wrapping", func(t *testing.T) {
		tape := buildFromReadmeTape([]string{`echo "hello"`}, opts)
		if !strings.Contains(tape, `Type 'echo "hello"'`) {
			t.Errorf("expected single-quoted Type, got:\n%s", tape)
		}
	})

	t.Run("multiple commands all appear", func(t *testing.T) {
		tape := buildFromReadmeTape([]string{"echo one", "echo two", "echo three"}, opts)
		for _, cmd := range []string{"echo one", "echo two", "echo three"} {
			if !strings.Contains(tape, cmd) {
				t.Errorf("tape missing command %q\ntape:\n%s", cmd, tape)
			}
		}
	})

	t.Run("custom pause appears in tape", func(t *testing.T) {
		custom := opts
		custom.pause = 3 * time.Second
		tape := buildFromReadmeTape([]string{"echo hi"}, custom)
		if !strings.Contains(tape, "Sleep 3s") {
			t.Errorf("expected Sleep 3s in tape:\n%s", tape)
		}
	})
}

// stringSliceEqual compares two string slices for equality, treating nil and
// empty slice as equal.
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
