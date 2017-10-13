package logger

import (
	"github.com/dgmann/ma-shared/stats"
)

type StatLogger struct {
	input <-chan stats.Stat
}

func NewStatLogger(input <-chan stats.Stat) StatLogger {
	return StatLogger{input:input}
}

func(logger *StatLogger) Print() {
	for stat := range logger.input {
		stat.Print()
	}
}

func(logger *StatLogger) Collect(print bool) *stats.Collector {
	collector := stats.NewCollector()
	for stat := range logger.input {
		if print {
			stat.Print()
		}
		collector.Add(stat)
	}
	return collector
}
