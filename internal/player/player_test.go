package player

import (
	"testing"
)

func TestNew(t *testing.T) {
	got := New()

	if got.x != 0 || got.y != 0 || got.movementCooldown != 0 {
		t.Errorf("Expected (x: 0, y: 0, movementCooldown: 0), got (x: %d, y: %d, movementCooldown: %d)",
			got.x, got.y, got.movementCooldown)
	}

	if len(got.walk) != 0 {
		t.Errorf("Expected walk to be nil or empty, got %+v", got.walk)
	}

	if got.Level != nil {
		t.Errorf("Expected Level to be nil, got %+v", got.Level)
	}
}
