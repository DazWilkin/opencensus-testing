package datadog

import (
	"log"
	"os"
	"time"

	"github.com/dazwilkin/opencensus/stats/view"
	datadog "gopkg.in/zorkian/go-datadog-api.v2"
)

var (
	client *datadog.Client
)

type Importer struct {
	name string
}

// NewImporter creates a new importer using the Options provided
func NewImporter(o Options) (*Importer, error) {
	return &Importer{
		name: "datadog",
	}, nil
}

// Name returns the Importer's name
func (i *Importer) Name() string {
	return i.name
}

// Value returns the Importer's value for the View, with the label values and the time specified
func (i *Importer) Value(v *view.View, labelValues []string, t time.Time) (float64, error) {

	from := t.Add(time.Minute * -1)

	query := NewQuery(v.Name)
	host, _ := os.Hostname()
	query.AddHostname(host)
	for i, labelName := range v.LabelNames {
		query.AddTagValue(labelName, labelValues[i])
	}
	log.Println(query.String())

	ss, err := client.QueryMetrics(from.Unix(), t.Unix(), query.String())
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range ss {
		log.Printf("Metric: %v", *s.Metric)
		for _, p := range s.Points {
			log.Printf("[%v] %v", time.Unix(0, int64(*p[0])*int64(time.Millisecond)), *p[1])
		}
	}

	return 0, nil
}

// Options represents the configuration of an OpenCensus Importer
type Options struct {
	MetricPrefix string
}

//TODO(dazwilkin) Instead of package init should this by a type func or helper?
func init() {
	API := os.Getenv("dd.API")
	App := os.Getenv("dd.App")
	if API == "" || App == "" {
		log.Fatal("Expect Datadog API and App values")
	}
	client = datadog.NewClient(API, App)
}
