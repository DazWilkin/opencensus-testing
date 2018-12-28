package view

import (
	"time"

	"github.com/golang/glog"
)

var (
	importers = make(map[string]Importer)
)

// Importer defines the interface for importers
type Importer interface {
	Name() string
	Value(v *View, labelValues []string, t time.Time) (float64, error)
}

// RegisterImporter adds an Importer to the View
func RegisterImporter(i Importer) {
	name := i.Name()
	if name == "" {
		glog.Fatal("Importer name must not be \"\"")
	}
	importers[i.Name()] = i
}

// UnregisterImporter removes an Importer from the View
func UnregisterImporter(i Importer) {
	name := i.Name()
	delete(importers, name)
}
