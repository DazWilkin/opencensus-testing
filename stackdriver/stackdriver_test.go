package stackdriver

import (
	"testing"
)

const (
	metricPrefix = "Freddie"
)

func Test_NewImporter(t *testing.T) {
	//TODO(dazwilkin) Add test(s) for returned importer
	t.Run("Null Options", func(t *testing.T) {
		i, _ := NewImporter(Options{})
		if got, want := i.name, "stackdriver"; got != want {
			t.Errorf("got %s; want %s", got, want)
		}
	})
	t.Run("With Options", func(t *testing.T) {
		i, _ := NewImporter(Options{
			MetricPrefix: metricPrefix,
		})
		if got, want := i.name, "stackdriver"; got != want {
			t.Errorf("got %s; want %s", got, want)
		}
		if got, want := i.options.MetricPrefix, metricPrefix; got != want {
			t.Errorf("got %s; want %s", got, want)
		}
	})
}
func TestView_Name(t *testing.T) {
	i, _ := NewImporter(Options{})
	if got, want := i.Name(), "stackdriver"; got != want {
		t.Errorf("got %s; want %s", got, want)
	}
}
func TestView_Value(t *testing.T) {

}
