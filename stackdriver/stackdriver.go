package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/golang/protobuf/ptypes/timestamp"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/api/iterator"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func tagsToFilter(labels, values []string) (string, error) {
	if len(labels) == 0 {
		return "", nil
	}
	if len(labels) != len(values) {
		return "", fmt.Errorf("Mismatched number of labels (%v) and values (%v)", len(labels), len(values))
	}
	res := ""
	for i, label := range labels {
		res += fmt.Sprintf("metric.label.%s=\"%s\"", label, values[i])
	}
	return res, nil
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
	filter := fmt.Sprintf(
		"resource.type=\"%s\" metric.type=\"%s\"",
		"global",
		"custom.googleapis.com/opencensus/181217_log_rpc_success_latency",
	)

	// This is how they're provided by the Interface
	labels := []string{"opencensus_task", "method"}
	values := []string{"go-117008@dazwilkin.kir.corp.google.com", "/trillian.TrillianAdmin/CreateTree"}
	labelvalues, err := tagsToFilter(labels, values)
	if err != nil {
		log.Fatalf("Unable to convert labels|values to filter string")
	}
	filter += labelvalues

	now := time.Now()
	fmt.Println(now.Format(time.RFC3339))
	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   name,
		Filter: filter,
		Interval: &monitoringpb.TimeInterval{
			EndTime: &googlepb.Timestamp{
				Seconds: now.Unix(),
			},
			StartTime: &googlepb.Timestamp{
				// One day ago
				Seconds: now.AddDate(0, 0, -1).Unix(),
			},
		},
	}
	it := client.ListTimeSeries(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp.Metric.Type)
		// Labels
		for key, value := range resp.Metric.Labels {
			fmt.Println("[" + key + "]=" + value)
		}
		// Points
		for _, point := range resp.Points {
			fmt.Printf("[%s:%s]=%s\n",
				intervalToString(point.Interval.StartTime),
				intervalToString(point.Interval.EndTime),
				//[TODO:dazwilkin] This is insufficient. Need to get its typedvalue ...
				point.Value,
			)
		}
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
