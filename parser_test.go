package dolly

import "testing"

func TestParser(t *testing.T) {
	input := `
Type "echo 'Hello, World!'" Enter
Right 3 Left 3
Sleep 3s`

	l := NewLexer(input)
	p := NewParser(l)

	cmds := p.Parse()
	errs := p.Errors()
	if len(errs) != 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}

	for _, cmd := range cmds {
		t.Log(cmd)
	}
}
