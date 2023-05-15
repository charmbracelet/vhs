package main

import (
	"reflect"
	"testing"
)

func TestCommand(t *testing.T) {
	const numberOfCommands = 24
	if len(CommandTypes) != numberOfCommands {
		t.Errorf("Expected %d commands, got %d", numberOfCommands, len(CommandTypes))
	}

	const numberOfCommandFuncs = 23
	if len(CommandFuncs) != numberOfCommandFuncs {
		t.Errorf("Expected %d commands, got %d", numberOfCommands, len(CommandTypes))
	}
}

func TestExecuteSetTheme(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		theme, err := getTheme("  ")
		requireNoErr(t, err)
		requireDefaultTheme(t, theme)
	})
	t.Run("named", func(t *testing.T) {
		theme, err := getTheme("Andromeda")
		requireNoErr(t, err)
		requireNotDefaultTheme(t, theme)
	})
	t.Run("json", func(t *testing.T) {
		theme, err := getTheme(`{"background": "#29283b"}`)
		requireNoErr(t, err)
		requireNotDefaultTheme(t, theme)
		if "#29283b" != theme.Background {
			t.Errorf("wrong background, expected %q, got %q", "#29283b", theme.Background)
		}
	})
	t.Run("suggestion", func(t *testing.T) {
		theme, err := getTheme("cattppuccin latt")
		requireEqualErr(t, err, "invalid `Set Theme \"cattppuccin latt\"`: did you mean \"Catppuccin Latte\"")
		requireDefaultTheme(t, theme)
	})
	t.Run("invalid json", func(t *testing.T) {
		theme, err := getTheme(`{"background`)
		requireErr(t, err)
		requireDefaultTheme(t, theme)
	})
	t.Run("unknown theme", func(t *testing.T) {
		theme, err := getTheme("foobar")
		requireErr(t, err)
		requireDefaultTheme(t, theme)
	})
}

func TestExecutePause(t *testing.T) {
	t.Run("Should create hiddenTerm if NOT exist and stop the execution", func(t *testing.T) {
		vhs := New()
		vhs.currentTerm = &Terminal{}

		ExecutePause(Command{Type: PAUSE}, &vhs)

		if vhs.executing {
			t.Errorf("Pause command has not stopped the execution")
		}

		if vhs.hiddenTerm == nil {
			t.Errorf("Pause command has not initialized hiddenTerm")
		}
	})
}

func TestExecuteResume(t *testing.T) {
	t.Run("Should use mainTerm as currentTerm and restart the execution", func(t *testing.T) {
		vhs := New()
		vhs.mainTerm = &Terminal{}

		ExecuteResume(Command{Type: RESUME}, &vhs)

		if !vhs.executing {
			t.Errorf("Resume command has not restared the execution")
		}

		if vhs.currentTerm != vhs.mainTerm {
			t.Errorf("Resume has not set mainTerm as currentTerm")
		}
	})
}

func requireErr(tb testing.TB, err error) {
	tb.Helper()
	if err == nil {
		tb.Fatalf("expected an error, got nil")
	}
}

func requireEqualErr(tb testing.TB, err1 error, err2 string) {
	tb.Helper()
	if err1 == nil {
		tb.Fatalf("expected an error, got nil")
	}
	if err1.Error() != err2 {
		tb.Fatalf("errors do not match: %q != %q", err1.Error(), err2)
	}
}

func requireNoErr(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("expected no error, got: %v", err)
	}
}

func requireDefaultTheme(tb testing.TB, theme Theme) {
	tb.Helper()
	if !reflect.DeepEqual(DefaultTheme, theme) {
		tb.Fatalf("expected theme to be the default theme, got something else: %+v", theme)
	}
}

func requireNotDefaultTheme(tb testing.TB, theme Theme) {
	tb.Helper()
	if reflect.DeepEqual(DefaultTheme, theme) {
		tb.Fatalf("expected theme to be different from the default theme, got the default instead")
	}
}
