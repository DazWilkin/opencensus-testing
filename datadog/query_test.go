package datadog

import (
	"strings"
	"testing"
)

func Test_AddHost(t *testing.T) {
	q := NewQuery("freddie")
	t.Run("Test Empty Host", func(t *testing.T) {
		host := ""
		q.AddHostname(host)
		if got, want := strings.Contains(q.TagString(), "host"), false; got != want {
			t.Errorf("got: %t; want: %t", got, want)
		}
	})
	t.Run("Test Non-Empty Host", func(t *testing.T) {
		host := "freddie"
		q.AddHostname(host)
		if got, want := strings.Contains(q.TagString(), "{host:"+host+"}"), true; got != want {
			t.Errorf("got: %t; want: %t", got, want)
		}

	})
}
func Test_AddTagValue(t *testing.T) {
	q := NewQuery("freddie")
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
