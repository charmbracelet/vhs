package main

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/vhs/parser"
)

func TestCommand(t *testing.T) {
	const numberOfCommands = 31
	if len(parser.CommandTypes) != numberOfCommands {
		t.Errorf("Expected %d commands, got %d", numberOfCommands, len(parser.CommandTypes))
	}

	const numberOfCommandFuncs = 31
	if len(CommandFuncs) != numberOfCommandFuncs {
		t.Errorf("Expected %d commands, got %d", numberOfCommandFuncs, len(CommandFuncs))
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

func TestExecuteSetPromptColor(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		v := New()
		err := ExecuteSetPromptColor(parser.Command{Args: "#ff8000"}, &v)
		requireNoErr(t, err)
		if v.Options.PromptColor != "#FF8000" {
			t.Errorf("wrong prompt color, expected %q, got %q", "#FF8000", v.Options.PromptColor)
		}
	})
	t.Run("without hash", func(t *testing.T) {
		v := New()
		err := ExecuteSetPromptColor(parser.Command{Args: "FF8000"}, &v)
		requireNoErr(t, err)
		if v.Options.PromptColor != "#FF8000" {
			t.Errorf("wrong prompt color, expected %q, got %q", "#FF8000", v.Options.PromptColor)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		v := New()
		err := ExecuteSetPromptColor(parser.Command{Args: "#5B56E"}, &v)
		requireErr(t, err)
		if v.Options.PromptColor != DefaultPromptColor {
			t.Errorf("prompt color should be unchanged, expected %q, got %q", DefaultPromptColor, v.Options.PromptColor)
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
