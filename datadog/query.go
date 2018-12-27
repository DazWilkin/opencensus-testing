package datadog

import "strings"

type Query struct {
	metric string
	tags   map[string]string
}

func NewQuery(metric string) *Query {
	tags := make(map[string]string)
	return &Query{
		metric: metric,
		tags:   tags,
	}
}
func (q *Query) AddHostname(host string) {
	// Add "host" as if it were another Tag
	q.AddTagValue("host", host)
}
func (q *Query) AddTagValue(tag, value string) {
	if tag != "" && value != "" {
		q.tags[tag] = value
	}
}
func (q *Query) TagString() string {
	// If no tags have been added, return "" not "{}"
	if len(q.tags) == 0 {
		return ""
	}
	// Otherwise append them to an array
	tags := make([]string, 0, len(q.tags))
	for key, value := range q.tags {
		tags = append(tags, key+":"+value)
	}
	// Then join the array elements, separated by "," and wrapped in "{...}"
	return "{" + strings.Join(tags, ",") + "}"
}
func (q *Query) String() string {
	return q.metric + q.TagString()
}
