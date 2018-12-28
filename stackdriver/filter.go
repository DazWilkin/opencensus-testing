package stackdriver

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
)

// Filter represents a Stackdriver Filter string
type Filter string

// NewFilter returns a new Filter (an empty string)
func NewFilter() *Filter {
	return new(Filter)
}

// add concatenates the string to the end of the Filter string
func (f *Filter) add(s string) {
	spacer := ""
	if !f.Empty() {
		spacer = " "
	}
	*f = Filter((string)(*f) + spacer + s)
}

// AddResourceType optionally adds a resource.type string to the Filter
func (f *Filter) AddResourceType(t string) {
	const (
		resourceType = "resource.type"
	)
	if strings.Contains(f.String(), resourceType) {
		glog.Fatalf("Stackdriver filters may only contain one '%s'", resourceType)
	}
	f.add(fmt.Sprintf("%s=\"%s\"", resourceType, t))
}

// AddMetricType optionally adds a metric.type corresponding to an OpenCensus custom metric to the Filter
func (f *Filter) AddMetricType(t string) {
	const (
		metricPath = "custom.googleapis.com/opencensus"
		metricType = "metric.type"
	)
	if strings.Contains(f.String(), metricType) {
		glog.Fatalf("Stackdriver filters may only contain one '%s'", metricType)
	}
	f.add(fmt.Sprintf("%s=\"%s/%s\"", metricType, metricPath, t))

}

// AddLabels optionally adds a set (as a map) of metric.label.[key]=[value] to the Filter
//TODO(dazwilkin) Check for duplicate metric.label.[key]
func (f *Filter) AddLabels(m map[string]string) {
	const (
		metricLabel = "metric.label"
	)
	labels := []string{}
	for label, value := range m {
		metricLabel := fmt.Sprintf("%s.\"%s\"=\"%s\"", metricLabel, label, value)
		labels = append(labels, metricLabel)
	}
	f.add(strings.Join(labels, " "))
}

// Empty returns true if the Filter is empty
func (f *Filter) Empty() bool {
	return len(f.String()) == 0
}

// String returns the Filter as a string
func (f *Filter) String() string {
	return (string)(*f)
}
