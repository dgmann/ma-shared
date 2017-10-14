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

func(logger *StatLogger) Print()  chan stats.Stat {
	output := make(chan stats.Stat, 10000)
	go func() {
		for stat := range logger.input {
			output <- stat
			stat.Print()
		}
		close(output)
	}()
	return output
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
