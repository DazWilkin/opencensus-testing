package datadog

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/dazwilkin/opencensus/stats/view"
	datadog "gopkg.in/zorkian/go-datadog-api.v2"
)

var (
	client *datadog.Client
)

type Importer struct {
	name    string
	options Options
}

// NewImporter creates a new importer using the Options provided
func NewImporter(o Options) (*Importer, error) {
	return &Importer{
		name:    "datadog",
		options: o,
	}, nil
}

// Name returns the Importer's name
func (i *Importer) Name() string {
	return i.name
}

// Stop closes the connection to the Datadog service
func (i *Importer) Stop() {
	log.Println("[Stop] Does nothing")
}

// Value returns the Importer's value for the View, with the label values and the time specified
func (i *Importer) Value(v *view.View, labelValues []string, t time.Time) (float64, error) {

	from := t.Add(time.Minute * -1)

	//TODO(dazwilkin) Datadog appears to append label names to the metric name, try it out
	query := NewQuery(func(v *view.View) string {
		// Name
		name := v.Name
		// Prefix Namespace, if one exists
		if i.options.Namespace != "" {
			name = i.options.Namespace + "." + name
		}
		// Append label names
		if len(v.LabelNames) > 0 {
			name = name + "_" + strings.Join(v.LabelNames, "_")
		}
		return name
	}(v))

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

	// If there is a time-series, grab the most recent one
	if len(ss) >= 1 {
		s := ss[0]
		log.Printf("Metric: %v", *s.Metric)
		// If the time-series contains any data points, grab the most recent one
		if len(s.Points) >= 1 {
			p := s.Points[0]
			// *p[0] == Unix epoch timestamp in ms
			// *p[1] == data
			log.Printf("[%v] %v", time.Unix(0, int64(*p[0])*int64(time.Millisecond)), *p[1])
			return *p[1], nil
		}
		return 0.0, nil
	}
	return 0.0, nil
}

// Options represents the configuration of an OpenCensus Importer
type Options struct {
	Namespace string
}

//TODO(dazwilkin) Instead of package init should this by a type func or helper?
func init() {
	API := os.Getenv("DD_API")
	App := os.Getenv("DD_APP")
	if API == "" || App == "" {
		log.Fatal("Expect Datadog API and App values")
	}
	client = datadog.NewClient(API, App)
}
