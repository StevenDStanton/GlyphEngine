package player

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	got := New()
	want := &Player{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Want %v got %v", want, got)
	}

}
