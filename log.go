package logs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/willf/pad"
)

const (
	// NewInstanceMsg sets the message to indicate the start of the log
	NewInstanceMsg = "START"
	// EndInstanceMsg sets the message to indicate the end of the log
	EndInstanceMsg = "END"
	// LogLevelDebug defines a normal debug log
	LogLevelDebug = "DEBUG"
	// LogLevelPanic defines a panic log
	LogLevelPanic = "PANIC"
	// LogLevelFatal defines a fatal log
	LogLevelFatal = "FATAL"
	// DateFormat defines the log date format
	DateFormat = time.RFC3339
)

var (
	osExit = os.Exit

	// Story  log
	Story Log
)

// Log represents information about a rest server log.
type Log struct {
	entries   []Entry
	folder    string
	Pad       int
	startTime int64
}

// Entry represents information about a rest server log entry.
type Entry struct {
	Level   string
	Message string
	Time    time.Time
}

func (l Log) getDate(t time.Time) string {
	return t.Format(DateFormat + " MST")
}

// Padding ...
func (l *Log) Padding(i int) {
	l.Pad = i
	return
}

// New creates new instance of Log
func New() *Log {
	var log Log

	log.entries = make([]Entry, 1)
	log.entries[0] = Entry{
		Level:   LogLevelDebug,
		Message: NewInstanceMsg,
		Time:    time.Now().UTC(),
	}

	log.Pad = 30
	log.startTime = log.TimeMs()

	return &log
}

func (l *Log) addEntry(level string, v ...interface{}) {
	l.entries = append(
		l.entries,
		Entry{
			Level:   level,
			Message: fmt.Sprint(v...),
			Time:    time.Now(),
		},
	)
}

// Entries returns all the entries
func (l *Log) Entries() []Entry {
	return l.entries
}

// Record ...
func (l *Log) Record(header string, v interface{}) error {
	l.Print(pad.Right(header, l.Pad, " "), v)
	return nil
}

// Print a regular log
func (l *Log) Print(v ...interface{}) {
	l.addEntry(LogLevelDebug, v...)
}

// Panic then throws a panic with the same message afterwards
func (l *Log) Panic(v ...interface{}) {
	l.addEntry(LogLevelPanic, v...)
	panic(fmt.Sprint(v...))
}

// LastEntry returns the last inserted log
func (l *Log) LastEntry() Entry {
	return l.entries[len(l.entries)-1]
}

// Count returns number of inserted logs
func (l *Log) Count() int {
	return len(l.entries)
}

// TimeMs Time in Micro Second
func (l *Log) TimeMs() int64 {
	return time.Now().UnixNano() / 1000000
}

// Dump will print all the messages to the io.
func (l *Log) Dump(printLogs ...bool) string {
	var (
		line, format string
		params       []interface{}
		messages     []string
		print        bool
	)

	if len(printLogs) > 0 {
		print = printLogs[0]
	} else {
		print = true
	}

	l.Print("Elapse ", l.TimeMs()-l.startTime, " ms")
	l.addEntry(LogLevelDebug, EndInstanceMsg+"\n")

	len := len(l.entries)
	for i := 0; i < len; i++ {
		format = "%s\t%s"
		params = []interface{}{
			l.getDate(l.entries[i].Time),
			l.entries[i].Level,
		}

		params = append(params, pad.Left(strconv.Itoa(i), 3, "0"))
		format = format + "  %s"

		params = append(params, l.entries[i].Message)

		format = format + "  %s\n"
		line = fmt.Sprintf(format, params...)

		if print {
			fmt.Print(line)
		}
		messages = append(messages, line)

	}

	l.entries = make([]Entry, 1)
	l.entries[0] = Entry{
		Level:   LogLevelDebug,
		Message: NewInstanceMsg,
		Time:    time.Now(),
	}

	l.startTime = l.TimeMs()

	return strings.Join(messages, "\n")
}
