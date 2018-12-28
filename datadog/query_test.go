package datadog

import (
	"strings"
	"testing"
)

const (
	hostName   = "host"
	metricName = "Freddie"
)

func Test_NewQuery(t *testing.T) {
	q := NewQuery(metricName)
	if got, want := q.metric, metricName; got != want {
		t.Errorf("got: %s; want: %s", got, want)
	}
}
func Test_AddHostname(t *testing.T) {
	q := NewQuery(metricName)
	t.Run("Test Empty Host", func(t *testing.T) {
		host := ""
		q.AddHostname(host)
		if got, want := strings.Contains(q.TagString(), hostName), false; got != want {
			t.Errorf("got: %t; want: %t", got, want)
		}
	})
	t.Run("Test Non-Empty Host", func(t *testing.T) {
		q.AddHostname(hostName)
		if got, want := strings.Contains(q.TagString(), "{host:"+hostName+"}"), true; got != want {
			t.Errorf("got: %t; want: %t", got, want)
		}

	})
}
func Test_AddTagValue(t *testing.T) {
	q := NewQuery(metricName)
	t.Run("Empty Tags", func(t *testing.T) {
		if got, want := q.TagString(), ""; got != want {
			t.Errorf("got: %s; want: %s", got, want)
		}
	})
	q.AddTagValue("X", "x")
	t.Run("Single Tag", func(t *testing.T) {
		if got, want := strings.Contains(q.TagString(), "X:x"), true; got != want {
			t.Errorf("got: %t; want: %t", got, want)
		}
		if got, want := strings.Contains(q.TagString(), ","), false; got != want {
			t.Errorf("got: %t; want: %t", got, want)
		}

	})
	q.AddTagValue("Y", "y")
	t.Run("Multiple Tags", func(t *testing.T) {
		if got, want := strings.Contains(q.TagString(), "X:x"), true; got != want {
			t.Errorf("got: %t; want: %t", got, want)
		}
		if got, want := strings.Contains(q.TagString(), "Y:y"), true; got != want {
			t.Errorf("got: %t; want: %t", got, want)
		}
		if got, want := strings.Contains(q.TagString(), ","), true; got != want {
			t.Errorf("got: %t; want: %t", got, want)
		}

	})

}
func Test_TagString(t *testing.T) {
	q := NewQuery(metricName)
	q.AddTagValue("X", "x")
	q.AddTagValue("Y", "y")
	if got, want := strings.Contains(q.TagString(), "X:x"), true; got != want {
		t.Errorf("got %t; want %t", got, want)
	}
	if got, want := strings.Contains(q.TagString(), "Y:y"), true; got != want {
		t.Errorf("got %t; want %t", got, want)
	}
}
func Test_String(t *testing.T) {

}
