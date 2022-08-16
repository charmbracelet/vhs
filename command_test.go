package dolly

import (
	"testing"
)

func TestCommand(t *testing.T) {
	if len(Commands) != 9 {
		t.Errorf("Expected 9 commands, got %d", len(Commands))
	}
}
