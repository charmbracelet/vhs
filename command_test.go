package dolly

import (
	"testing"
)

func TestCommand(t *testing.T) {
	if len(Commands) != 10 {
		t.Errorf("Expected 10 commands, got %d", len(Commands))
	}
}
