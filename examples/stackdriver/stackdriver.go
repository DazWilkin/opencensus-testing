package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	importer_stackdriver "github.com/dazwilkin/opencensus/stackdriver"
	importer_view "github.com/dazwilkin/opencensus/stats/view"
	"github.com/golang/glog"
)

const (
	metricPrefix = "namespace"
)

func main() {

	// Stackdriver Exporter
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		MetricPrefix: metricPrefix,
	})
	if err != nil {
		glog.Fatal(err)
	}
	// Important to invoke Flush before exiting
	defer exporter.Flush()

	view.RegisterExporter(exporter)

	// Stackdriver requires 60s reporting period
	view.SetReportingPeriod(60 * time.Second)

	// Stackdriver Importer
	// Stackdriver metric.type == custom.googleapis.com/opencensus/181220_counter0
	// No reference to MetricPrefix
	importer, err := importer_stackdriver.NewImporter(importer_stackdriver.Options{
		MetricPrefix: metricPrefix,
	})
	if err != nil {
		glog.Fatal(err)
	}

	importer_view.RegisterImporter(importer)

	name := "counter0"
	measure := stats.Float64(name, "Testing", "1")

	labelNames := []string{"key1", "key2"}
	labelValues := []string{"value1", "value2"}

	tagKeys := []tag.Key{}
	for _, labelName := range labelNames {
		tagKey, err := tag.NewKey(labelName)
		if err != nil {
			glog.Fatal(err)
		}
		tagKeys = append(tagKeys, tagKey)
	}

	v := &view.View{
		Name:        name,
		Measure:     measure,
		Description: "Testing",
		Aggregation: view.Sum(),
		TagKeys:     tagKeys,
	}
	if err := view.Register(v); err != nil {
		glog.Fatal(err)
	}

	iv := &importer_view.View{
		Name:       name,
		LabelNames: labelNames,
	}
	if err := importer_view.Register(iv); err != nil {
		glog.Fatal(err)
	}

	ctx := context.TODO()
	for i, tagKey := range tagKeys {
		ctx, err = tag.New(ctx, tag.Insert(tagKey, labelValues[i]))
		if err != nil {
			glog.Fatal(err)
		}
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var wg sync.WaitGroup
	wg.Add(2)

	end := time.Now().Add(time.Minute * 10)

	// Measure
	go func(end time.Time) {
		defer wg.Done()
		sum := 0.0
		for end.After(time.Now()) {
			val := r1.Float64()
			sum += val
			glog.Infof("write: %f [%f]", val, sum)
			stats.Record(ctx, measure.M(val))
			// Write measurements every 10 seconds
			time.Sleep(10 * time.Second)
		}
		glog.Infof("Done Measuring")
	}(end)

	// Value
	go func(end time.Time) {
		defer wg.Done()
		for end.After(time.Now()) {
			// Read values after every 30 seconds
			time.Sleep(30 * time.Second)
			val, err := importer.Value(
				iv,
				labelValues,
				time.Now(),
			)
			message := ""
			if err != nil {
				message = fmt.Sprintf(" [%s]", err)
			}
			glog.Infof("reads: %f%s", val, message)
		}
		glog.Infof("Done Reading")
	}(end)

	wg.Wait()
}
