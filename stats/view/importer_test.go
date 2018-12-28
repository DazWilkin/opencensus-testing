package view

import (
	"testing"
	"time"
)

// Importer implements the Importer interface in order to be able to test the interface
type importer struct {
	name string
}

func (i *importer) Name() string {
	return i.name
}
func (i *importer) Value(v *View, labelValues []string, t time.Time) (float64, error) {
	return 0.0, nil
}
func Test_RegisterImporter(t *testing.T) {
	var i *importer
	t.Run("Empty Name", func(t *testing.T) {
		i = &importer{
			name: "",
		}
		RegisterImporter(i)
		_, got := importers[""]
		want := false
		if got != want {
			t.Errorf("got %t; want %t", got, want)
		}
	})
	t.Run("Acceptable Name", func(t *testing.T) {
		i = &importer{
			name: "X",
		}
		RegisterImporter(i)
		_, got := importers["X"]
		want := true
		if got != want {
			t.Errorf("got %t; want %t", got, want)
		}
	})
}
func Test_UnregisterImporter(t *testing.T) {
	var i *importer
	t.Run("Acceptable Name", func(t *testing.T) {
		i = &importer{
			name: "X",
		}
		RegisterImporter(i)
		UnregisterImporter(i)
		_, got := importers["X"]
		want := false
		if got != want {
			t.Errorf("got %t; want %t", got, want)
		}

	})
}
