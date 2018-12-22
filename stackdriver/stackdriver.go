package stackdriver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/dazwilkin/opencensus/stats/view"
	"github.com/golang/glog"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/api/metric"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

var (
	client    *monitoring.MetricClient
	projectID string
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

// Importer represents the inverse of an OpenCensus Exporter
// It gets values for measurements from the service
// For Stackdriver, we'll use ADCs but need a robot with >= Monitoring Viewer
type Importer struct {
	name string
}

// NewImporter creates a new importer using the Options provided
func NewImporter(o Options) (*Importer, error) {
	return &Importer{
		name: "stackdriver",
	}, nil
}

// Name returns the Importer's name
func (i *Importer) Name() string {
	return i.name
}

// Value returns the Importer's value for the View, with the label values and the time specified
func (i *Importer) Value(v *view.View, labelValues []string, t time.Time) (float64, error) {
	// Private functions
	createInterval := func(start, end time.Time) *monitoringpb.TimeInterval {
		return &monitoringpb.TimeInterval{
			StartTime: &googlepb.Timestamp{
				Seconds: start.Unix(),
			},
			EndTime: &googlepb.Timestamp{
				Seconds: end.Unix(),
			},
		}
	}
	getFloat64Value := func(t metric.MetricDescriptor_ValueType, p *monitoringpb.Point) (float64, error) {
		switch t {
		case metricpb.MetricDescriptor_DISTRIBUTION:
			dist := p.GetValue().GetDistributionValue()
			count := dist.GetCount()
			mean := dist.GetMean()
			return float64(count) * mean, nil
		case metricpb.MetricDescriptor_DOUBLE:
			return p.GetValue().GetDoubleValue(), nil
		case metricpb.MetricDescriptor_INT64:
			return float64(p.GetValue().GetInt64Value()), nil
		default:
			//TODO(dazwilkin) There are more types to enumerate
			return 0.0, nil
		}
	}
	mapLabelsValues := func(labels, values []string) map[string]string {
		m := map[string]string{}
		// Only proceed if there
		// - are labels and values to map
		// - is no discrepancy between the set of labels and values
		if labels == nil && values == nil {
			return m
		}
		if len(labels) != len(values) {
			glog.Fatal("Inconsistency between labels and values")
		}
		for i, label := range labels {
			m[label] = values[i]
		}
		return m
	}

	f := NewFilter()
	f.AddResourceType("global")
	f.AddMetricType(v.Name)

	// Convert Labels[],Values[]-->map(Label=Value)
	f.AddLabels(mapLabelsValues(v.LabelNames, labelValues))

	fmt.Println(t.Format(time.RFC3339))
	req := &monitoringpb.ListTimeSeriesRequest{
		Name:     fmt.Sprintf("projects/%s", projectID),
		Filter:   f.String(),
		Interval: createInterval(t.Add(time.Minute*-1), t),
	}
	it := client.ListTimeSeries(context.TODO(), req)

	// We only want the most-recent entry in the timeseries
	resp, err := it.Next()
	if err == iterator.Done {
		// There are no results
		return 0.0, errors.New("No timeseries match the filter")
	}
	if err != nil {
		// Something untoward
		return 0.0, err
	}

	// And only the most recent point from the most recent entry
	return getFloat64Value(resp.GetValueType(), resp.Points[0])
}

// Options represents the configuration of an OpenCensus Importer
type Options struct {
	MetricPrefix string
}

//TODO(dazwilkin) Instead of package init should this by a type func or helper?
func init() {
	ctx := context.Background()

	var err error
	// Assumes Application Default Credentials
	// Commonly credentials are provided using environment variable GOOGLE_APPLICATION_CREDENTIALS
	client, err = monitoring.NewMetricClient(ctx)
	if err != nil {
		glog.Fatal(err)
	}

	projectID = os.Getenv("PROJECT")
	if projectID == "" {
		glog.Fatal("Google Cloud Project ID is required; specify using environment variable 'PROJECT'")
	}
}
