package main

import (
	"testing"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/vt"
)

func TestCtrlKey(t *testing.T) {
	tests := []struct {
		args string
		want uv.KeyPressEvent
		ok   bool
	}{
		{"c", uv.KeyPressEvent{Code: 'c', Mod: uv.ModCtrl}, true},
		{"C", uv.KeyPressEvent{Code: 'c', Mod: uv.ModCtrl | uv.ModShift}, true},
		{"Shift Alt C", uv.KeyPressEvent{Code: 'c', Mod: uv.ModCtrl | uv.ModShift | uv.ModAlt}, true},
		{"Alt Left", uv.KeyPressEvent{Code: vt.KeyLeft, Mod: uv.ModCtrl | uv.ModAlt}, true},
		{"Enter", uv.KeyPressEvent{Code: vt.KeyEnter, Mod: uv.ModCtrl}, true},
		{"Space", uv.KeyPressEvent{Code: vt.KeySpace, Mod: uv.ModCtrl}, true},
		{"Backspace", uv.KeyPressEvent{Code: vt.KeyBackspace, Mod: uv.ModCtrl}, true},
		{"@", uv.KeyPressEvent{Code: vt.KeySpace, Mod: uv.ModCtrl}, true},
		{"-", uv.KeyPressEvent{Code: '_', Mod: uv.ModCtrl}, true},
		{"[", uv.KeyPressEvent{Code: '[', Mod: uv.ModCtrl}, true},
		{"Shift", uv.KeyPressEvent{Mod: uv.ModCtrl}, false},
		{"", uv.KeyPressEvent{Mod: uv.ModCtrl}, false},
	}
	for _, tt := range tests {
		got, ok := ctrlKey(tt.args)
		if got != tt.want || ok != tt.ok {
			t.Errorf("%q: expected %v (%t), got %v (%t)", tt.args, tt.want, tt.ok, got, ok)
		}
	}
}

func TestRuneKey(t *testing.T) {
	tests := []struct {
		r    rune
		mod  uv.KeyMod
		want uv.KeyPressEvent
	}{
		{'a', 0, uv.KeyPressEvent{Code: 'a', Text: "a"}},
		{'a', uv.ModShift, uv.KeyPressEvent{Code: 'A', Text: "A", Mod: uv.ModShift}},
		{'1', uv.ModShift, uv.KeyPressEvent{Code: '!', Text: "!", Mod: uv.ModShift}},
		{'/', uv.ModShift, uv.KeyPressEvent{Code: '?', Text: "?", Mod: uv.ModShift}},
		{'.', uv.ModAlt, uv.KeyPressEvent{Code: '.', Text: ".", Mod: uv.ModAlt}},
		{'é', 0, uv.KeyPressEvent{Code: 'é', Text: "é"}},
		{'\n', 0, keyEnter},
		{'\t', uv.ModShift, uv.KeyPressEvent{Code: vt.KeyTab, Mod: uv.ModShift}},
		{'↑', 0, keyUp},
	}
	for _, tt := range tests {
		if got := runeKey(tt.r, tt.mod); got != tt.want {
			t.Errorf("%q: expected %v, got %v", tt.r, tt.want, got)
		}
	}
}
