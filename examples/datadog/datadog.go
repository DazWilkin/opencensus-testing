package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	exporter_datadog "github.com/DataDog/opencensus-go-exporter-datadog"
	importer_datadog "github.com/dazwilkin/opencensus/datadog"
	importer_view "github.com/dazwilkin/opencensus/stats/view"
	"github.com/golang/glog"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

const (
	namespace = "namespace"
)

func main() {

	// Datadog Exporter
	exporter, err := exporter_datadog.NewExporter(exporter_datadog.Options{
		Namespace: namespace,
	})
	if err != nil {
		glog.Fatal(err)
	}
	defer exporter.Stop()

	view.RegisterExporter(exporter)

	// Datadog Importer
	importer, err := importer_datadog.NewImporter(importer_datadog.Options{
		Namespace: namespace,
	})
	if err != nil {
		glog.Fatal(err)
	}
	defer importer.Stop()

	name := "counter0"
	measure := stats.Float64(name, "Testing", "1")

	labelNames := []string{"key1", "key2"}
	labelValues := []string{"value1", "value2"}

	tagKeys := []tag.Key{}
	for _, label := range labelNames {
		key, err := tag.NewKey(label)
		if err != nil {
			glog.Fatal(err)
		}
		tagKeys = append(tagKeys, key)
	}

	// Exporter's View
	ev := &view.View{
		Name:        name,
		Measure:     measure,
		Description: "Testing",
		Aggregation: view.Sum(),
		TagKeys:     tagKeys,
	}
	if err := view.Register(ev); err != nil {
		glog.Fatal(err)
	}

	// Importer's View
	iv := &importer_view.View{
		Name:       name,
		LabelNames: labelNames,
	}
	if err := importer_view.Register(iv); err != nil {
		glog.Fatal(err)
	}

	ctx := context.TODO()
	for i, key := range tagKeys {
		ctx, err = tag.New(ctx, tag.Insert(key, labelValues[i]))
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
