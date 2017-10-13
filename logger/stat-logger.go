package logger

import (
	"github.com/dgmann/ma-shared"
	"fmt"
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

func(logger *StatLogger) Collect() *stats.Collector {
	collector := stats.NewCollector()
	go func() {
		for stat := range logger.input {
			collector.Add(stat)
		}
	}()
	return collector
}
