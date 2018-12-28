package view

import (
	"testing"
)

func Test_Register(t *testing.T) {
	var v *View
	t.Run("Empty Name", func(t *testing.T) {
		v = &View{
			Name: "",
		}
		if got, want := Register(v), "View name must not be \"\""; got.Error() != want {
			t.Errorf("got %s, wanted %s", got, want)
		}
	})
	t.Run("Acceptable Name", func(t *testing.T) {
		v = &View{
			Name: "X",
		}
		Register(v)
		if got, want := views["X"].Name, "X"; got != want {
			t.Errorf("got %s, wanted %s", got, want)
		}
	})
}
