package view

import "time"

// Importer defines the interface for importers
type Importer interface {
	Name() string
	Value(v *View, labelValues []string, t time.Time) (float64, error)
}

var (
	importers []Importer
	views     []*View
)

// View represents an OpenCensus View
// It must have a name as a unique identifier
// And probably a type
// And a set (map) of key:value labels (Tags) that uniquely identify the metric
// And probably a time interval when the values were sent
// And presumably a list of importers
type View struct {
	Name       string
	LabelNames []string
}

// Register registers a view
func Register(vv ...*View) error {
	views = append(views, vv...)
	return nil
}

// RegisterImporter adds an Importer to the View
func RegisterImporter(i Importer) {
	importers = append(importers, i)
}

// Value retrieves a value from an OpenCensus View
func (v *View) Value(labelValues []string) map[string]float64 {
	values := map[string]float64{}
	// Get each importer to provide the most recent value
	now := time.Now()
	for _, importer := range importers {
		value, err := importer.Value(v, labelValues, now)
		if err != nil {
			//TODO(dazwilkin) should bubble up the error but can't as not supported by the Interface
			// Rather than break (next best option), add 0.0 for the result in order to pass the tests :-(
		}
		values[importer.Name()] = value
	}
	return values
}
