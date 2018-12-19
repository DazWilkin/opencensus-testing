package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/golang/protobuf/ptypes/timestamp"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/api/iterator"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// Filter represents a Stackdriver filter string
//TODO(dazwilkin) Would it be preferable to represent Stackdriver (string) filters as a type (with methods)?
type filter string

// newFilter returns a new filter (an empty string)
func newFilter() *filter {
	return new(filter)
}

// addResourceType optionally adds a resource.type string to the filter
func (f *filter) addResourceType(t string) {
	*f = filter(fmt.Sprintf("%s resource.type=\"%s\"", *f, t))
}

// addMetricType optionally adds a metric.type corresponding to an OpenCensus custom metric to the filter
func (f *filter) addMetricType(t string) {
	const (
		metricPath = "custom.googleapis.com/opencensus"
	)
	*f = filter(fmt.Sprintf("%s %s/%s", *f, metricPath, t))

}

// addLabels optionally adds a set (as a map) of metric.label.[key]=[value] to the filter
func (f *filter) addLabels(m map[string]string) {
	labels := []string{}
	for label, value := range m {
		metricLabel := fmt.Sprintf("metric.label.%s=\"%s\"", label, value)
		labels = append(labels, metricLabel)
	}
	*f = *f + filter(strings.Join(labels, " "))
}

// String returns the filter as a string
func (f *filter) String() string {
	return (string)(*f)
}

// View represents an OpenCensus View
// It must have a name as a unique identifier
// And probably a type
// And a set (map) of key:value labels (Tags) that uniquely identify the metric
// And probably a time interval when the values were sent
type View struct {
	name string
	tags map[string]string
}

// Value retrieves a value from an OpenCensus View
func (v *View) Value() float64 {
	return 0.0
}

func metricNameToFilter(prefix, name string) string {
	const (
		resourceType = "global"
		metricPath   = "custom.googleapis.com/opencensus"
	)
	return fmt.Sprintf("resource.type=\"%s\" metric.type=\"%s\"",
		resourceType,
		fmt.Sprintf("%s/%s%s", metricPath, prefix, name),
	)
}
func tagsToFilter(labels, values []string) (string, error) {
	if len(labels) == 0 {
		return "", nil
	}
	if len(labels) != len(values) {
		return "", fmt.Errorf("Mismatched number of labels (%v) and values (%v)", len(labels), len(values))
	}
	res := ""
	for i, label := range labels {
		// The prefixing space is important to ensure separation of values in the filter string
		res += fmt.Sprintf(" metric.label.%s=\"%s\"", label, values[i])
	}
	return res, nil
}
func createInterval(start, end time.Time) *monitoringpb.TimeInterval {
	return &monitoringpb.TimeInterval{
		StartTime: &googlepb.Timestamp{
			// One minute ago
			Seconds: start.Unix(),
		},
		EndTime: &googlepb.Timestamp{
			Seconds: end.Unix(),
		},
	}
}
func intervalToString(ts *timestamp.Timestamp) string {
	return time.Unix(ts.Seconds, 0).Format(time.RFC3339)
}

func main() {

	projectID := os.Getenv("PROJECT")

	ctx := context.Background()

	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	name := fmt.Sprintf("projects/%s", projectID)

	prefix := "181217_"

	metrics := []struct {
		name   string
		labels []string
		values []string
	}{
		{
			name:   "log_rpc_requests",
			labels: []string{"opencensus_task", "method"},
			// Differs in method key's value
			values: []string{"go-125577@dazwilkin.kir.corp.google.com", "/trillian.TrillianAdmin/CreateTree"},
		},
		{
			name:   "log_rpc_requests",
			labels: []string{"opencensus_task", "method"},
			// Differs in method key's value
			values: []string{"go-125577@dazwilkin.kir.corp.google.com", "/trillian.TrillianLog/QueueLeaf"},
		},
		{
			name:   "log_rpc_success_latency",
			labels: []string{"opencensus_task", "method"},
			values: []string{"go-117008@dazwilkin.kir.corp.google.com", "/trillian.TrillianAdmin/CreateTree"},
		},
	}
	for _, metric := range metrics {
		filter := metricNameToFilter(prefix, metric.name)
		// This is how they're provided by the Interface

		labelvalues, err := tagsToFilter(metric.labels, metric.values)
		if err != nil {
			log.Fatalf("Unable to convert labels|values to filter string")
		}
		filter += labelvalues
		now := time.Now()
		fmt.Println(now.Format(time.RFC3339))
		req := &monitoringpb.ListTimeSeriesRequest{
			Name:     name,
			Filter:   filter,
			Interval: createInterval(now.Add(time.Hour*-36), now),
		}
		it := client.ListTimeSeries(ctx, req)
		// Only interested in the most-recent result
		// for {
		resp, err := it.Next()
		if err == iterator.Done {
			// There are no results
			// break
			log.Println("No results")
		}
		if err != nil {
			log.Println(err)
		}
		log.Println(resp.Metric.Type)
		// Labels
		for key, value := range resp.Metric.Labels {
			fmt.Println("[" + key + "]=" + value)
		}
		var processPoint func(*monitoringpb.Point) string
		if resp.GetValueType() == metricpb.MetricDescriptor_DISTRIBUTION {
			processPoint = func(p *monitoringpb.Point) string {
				dist := p.GetValue().GetDistributionValue()
				count := dist.GetCount()
				mean := dist.GetMean()
				return fmt.Sprintf("[%s:%s] count: %d, mean: %f, sum: %f",
					intervalToString(p.Interval.StartTime),
					intervalToString(p.Interval.EndTime),
					count,
					mean,
					float64(count)*mean,
				)
			}
		} else {
			processPoint = func(p *monitoringpb.Point) string {
				return fmt.Sprintf("[%s:%s]=%s",
					intervalToString(p.Interval.StartTime),
					intervalToString(p.Interval.EndTime),
					//[TODO:dazwilkin] This is insufficient. Need to get its typedvalue ...
					p.Value,
				)
			}
		}

		for _, point := range resp.Points {
			fmt.Println(processPoint(point))
		}

		// }
	}
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
