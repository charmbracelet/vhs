package main

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/vhs/parser"
	"github.com/charmbracelet/vhs/token"
)

func TestCommand(t *testing.T) {
	const numberOfCommands = 32
	if len(parser.CommandTypes) != numberOfCommands {
		t.Errorf("Expected %d commands, got %d", numberOfCommands, len(parser.CommandTypes))
	}

	const numberOfCommandFuncs = 32
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

func TestExecuteSetRowsColumns(t *testing.T) {
	t.Run("rows", func(t *testing.T) {
		v := New()
		err := ExecuteSetRows(parser.Command{Type: token.SET, Options: "Rows", Args: "40"}, &v)
		requireNoErr(t, err)
		if v.Options.Video.Style.Rows != 40 {
			t.Errorf("expected Rows to be 40, got %d", v.Options.Video.Style.Rows)
		}
	})

	t.Run("columns", func(t *testing.T) {
		v := New()
		err := ExecuteSetColumns(parser.Command{Type: token.SET, Options: "Columns", Args: "100"}, &v)
		requireNoErr(t, err)
		if v.Options.Video.Style.Columns != 100 {
			t.Errorf("expected Columns to be 100, got %d", v.Options.Video.Style.Columns)
		}
	})

	t.Run("rows must be positive", func(t *testing.T) {
		v := New()
		err := ExecuteSetRows(parser.Command{Type: token.SET, Options: "Rows", Args: "0"}, &v)
		requireErr(t, err)
	})

	t.Run("columns must be positive", func(t *testing.T) {
		v := New()
		err := ExecuteSetColumns(parser.Command{Type: token.SET, Options: "Columns", Args: "0"}, &v)
		requireErr(t, err)
	})

	t.Run("height then rows conflict", func(t *testing.T) {
		v := New()
		requireNoErr(t, ExecuteSetHeight(parser.Command{Type: token.SET, Options: "Height", Args: "600"}, &v))
		err := ExecuteSetRows(parser.Command{Type: token.SET, Options: "Rows", Args: "40"}, &v)
		requireErr(t, err)
		if v.Options.Video.Style.Rows != 0 {
			t.Errorf("expected Rows to remain unset, got %d", v.Options.Video.Style.Rows)
		}
	})

	t.Run("rows then height conflict", func(t *testing.T) {
		v := New()
		requireNoErr(t, ExecuteSetRows(parser.Command{Type: token.SET, Options: "Rows", Args: "40"}, &v))
		defaultHeight := v.Options.Video.Style.Height
		err := ExecuteSetHeight(parser.Command{Type: token.SET, Options: "Height", Args: "700"}, &v)
		requireErr(t, err)
		if v.Options.Video.Style.Height != defaultHeight {
			t.Errorf("expected Height to remain unchanged, got %d", v.Options.Video.Style.Height)
		}
	})

	t.Run("width then columns conflict", func(t *testing.T) {
		v := New()
		requireNoErr(t, ExecuteSetWidth(parser.Command{Type: token.SET, Options: "Width", Args: "1200"}, &v))
		err := ExecuteSetColumns(parser.Command{Type: token.SET, Options: "Columns", Args: "100"}, &v)
		requireErr(t, err)
		if v.Options.Video.Style.Columns != 0 {
			t.Errorf("expected Columns to remain unset, got %d", v.Options.Video.Style.Columns)
		}
	})

	t.Run("columns then width conflict", func(t *testing.T) {
		v := New()
		requireNoErr(t, ExecuteSetColumns(parser.Command{Type: token.SET, Options: "Columns", Args: "100"}, &v))
		defaultWidth := v.Options.Video.Style.Width
		err := ExecuteSetWidth(parser.Command{Type: token.SET, Options: "Width", Args: "1300"}, &v)
		requireErr(t, err)
		if v.Options.Video.Style.Width != defaultWidth {
			t.Errorf("expected Width to remain unchanged, got %d", v.Options.Video.Style.Width)
		}
	})

	t.Run("width and rows can coexist", func(t *testing.T) {
		v := New()
		requireNoErr(t, ExecuteSetWidth(parser.Command{Type: token.SET, Options: "Width", Args: "1200"}, &v))
		requireNoErr(t, ExecuteSetRows(parser.Command{Type: token.SET, Options: "Rows", Args: "40"}, &v))
	})

	t.Run("height and columns can coexist", func(t *testing.T) {
		v := New()
		requireNoErr(t, ExecuteSetHeight(parser.Command{Type: token.SET, Options: "Height", Args: "600"}, &v))
		requireNoErr(t, ExecuteSetColumns(parser.Command{Type: token.SET, Options: "Columns", Args: "100"}, &v))
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
