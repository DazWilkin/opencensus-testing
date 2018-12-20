package stackdriver

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewFilter(t *testing.T) {
	f := NewFilter()
	if got, want := f.String(), ""; got != want {
		t.Errorf("[NewFilter] got=\"%s\" want=\"%s\"", got, want)
	}
}
func TestFilter_add(t *testing.T) {
	f := NewFilter()
	t.Run("add(\"\")", func(t *testing.T) {
		// Should leave Filter unchanged
		f.add("")
		if got, want := f.String(), ""; got != want {
			t.Errorf("[add] got=\"%s\" want=\"%s\"", got, want)
		}
	})
	t.Run("add(\"X\")", func(t *testing.T) {
		f.add("X")
		if got, want := f.String(), "X"; got != want {
			t.Errorf("[add] got=\"%s\" want=\"%s\"", got, want)
		}
	})
	t.Run("add(\"Y\")", func(t *testing.T) {
		// Tests the addition of a spacer when the string is non-empty
		f.add("Y")
		if got, want := f.String(), "X Y"; got != want {
			t.Errorf("[add] got=\"%s\" want=\"%s\"", got, want)
		}
	})
}
func TestFilter_AddResourceType(t *testing.T) {
	f := NewFilter()
	f.AddResourceType("X")
	if got, want := f.String(), "resource.type=\"X\""; got != want {
		t.Errorf("[addResourceType] got=\"%s\" want=\"%s\"", got, want)
	}
}
func TestFilter_AddMetricType(t *testing.T) {
	f := NewFilter()
	f.AddMetricType("X")
	if got, want := f.String(), "metric.type=\"custom.googleapis.com/opencensus/X\""; got != want {
		t.Errorf("[addMetricType] got=\"%s\" want=\"%s\"", got, want)
	}
}
func TestFilter_AddLabels(t *testing.T) {
	f := NewFilter()

	t.Run("No Labels", func(t *testing.T) {
		m := map[string]string{}
		// Empty set of labels should leave Filter unchanged
		f.AddLabels(m)
		if got, want := f.String(), ""; got != want {
			t.Errorf("[addLabels] got=\"%s\" want=\"%s\"", got, want)
		}
	})
	t.Run("Some Labels", func(t *testing.T) {
		// Non-zero set of labels should be add to the Filter (in order although ordering is unimportant)
		m := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		f.AddLabels(m)
		for key, value := range m {
			// Ordering is not relevant so using 'Contains' to find string existence
			if got, want := strings.Contains(
				f.String(),
				fmt.Sprintf("metric.label.%s=\"%s\"", key, value),
			), true; got != want {
				t.Errorf("[addLabels] Unable to find 'metric.label.%s=\"%s\"'", key, value)
			}
		}
	})
}
func TestFilter_Empty(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		f := NewFilter()
		if got, want := f.Empty(), true; got != want {
			t.Errorf("[Empty] Expected Filter to be empty")
		}
	})
	t.Run("Non-Empty", func(t *testing.T) {
		f := NewFilter()
		f.add("X")
		if got, want := f.Empty(), false; got != want {
			t.Errorf("[Empty] Expected Filter to not be empty")
		}

	})
}
func TestFilter_String(t *testing.T) {
	f := NewFilter()
	t.Run("Empty Filter", func(t *testing.T) {
		if got, want := f.String(), ""; got != want {
			t.Errorf("[String] Expected an empty string")
		}
	})
	t.Run("Non-Empty Filter", func(t *testing.T) {
		f.add("X")
		if got, want := f.String(), "X"; got != want {
			t.Errorf("[String] Expected an empty string")
		}
	})
	t.Run("Non-Empty Filter w/ Spacer", func(t *testing.T) {
		f.add("Y")
		if got, want := f.String(), "X Y"; got != want {
			t.Errorf("[String] Expected an empty string")
		}
	})

}
func TestFilter_AddMultipleParts(t *testing.T) {
	f := NewFilter()
	f.AddResourceType("R")
	f.AddMetricType("M")
	f.AddLabels(map[string]string{
		"key1": "value1",
		"key2": "value2",
	})
	if got, want := f.String(), "resource.type=\"R\" metric.type=\"custom.googleapis.com/opencensus/M\" metric.label.key1=\"value1\" metric.label.key2=\"value2\""; got != want {
		t.Errorf("[addMultipleParts] got=\"%s\" want=\"%s\"", got, want)
	}
}

func TestView_Value(t *testing.T) {

}

func TestNewImporter(t *testing.T) {
	//TODO(dazwilkin) Add test(s) for returned importer
	_, err := NewImporter(Options{})
	if err != nil {
		t.Errorf("[NewImporter] Return error unexpectedly")
	}

}
