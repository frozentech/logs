package logs

import (
	"fmt"
	"strings"
	"time"
)

// Entry ...
type Entry struct {
	Name    string      `json:"name"`
	Message interface{} `json:"message"`
	Time    time.Time   `json:"timestamp"`
}

// Stories ...
var Stories *Logs

// Record ...
func Record(name string, v interface{}) {
	if Stories != nil {
		Stories.Push(name, v)
	}
}

// Logs ...
type Logs struct {
	Entries []Entry
}

// New ...
func New() *Logs {
	var log Logs
	log.Entries = make([]Entry, 0)
	return &log
}

// Push ...
func (l *Logs) Push(name string, v interface{}) {
	l.Entries = append(
		l.Entries,
		Entry{
			Name:    name,
			Message: v,
			Time:    time.Now().UTC(),
		},
	)
}

func format(name string) string {
	return strings.Replace(name, " ", "_", -1)
}

// Dump ...
func (l *Logs) Dump() (fields map[string]interface{}) {
	fields = make(map[string]interface{}, 0)

	for idx, item := range l.Entries {
		if item.Message == nil {
			fields[fmt.Sprintf("%06d", idx)] = item.Name
		} else {
			fields[fmt.Sprintf("%06d.%s", idx, item.Name)] = item.Message
		}
	}

	return
}
