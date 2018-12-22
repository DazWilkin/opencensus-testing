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
func Test_Register(t *testing.T)         {}
func Test_RegisterImporter(t *testing.T) {}
