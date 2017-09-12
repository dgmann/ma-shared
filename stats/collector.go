package stats

import (
	"encoding/csv"
	"os"
	"log"
	"fmt"
	"github.com/dgmann/ma-shared"
	"io"
)

type Collector struct {
	Stats []Stat
}

func NewCollector() *Collector {
	return &Collector{Stats: make([]Stat, 0, 1500)}
}

func(collector *Collector) Add(stat Stat) {
	collector.Stats = append(collector.Stats, stat)
}

func(collector *Collector) ToCSV(w io.Writer) {
	writer := csv.NewWriter(w)

	if len(collector.Stats) > 0 {
		stat := collector.Stats[0]
		header := make([]string, 0, len(stat.Stages))
		for _, stage := range stat.Stages {
			header = append(header, stage.Name)
		}
		header = append(header, "Total")
		if err := writer.Write(header); err != nil {
			log.Fatalln("error writing header to csv:", err)
		}
	}
	for _, stat := range collector.Stats {
		record := []string{}
		for _, stage := range stat.Stages {
			record = append(record, fmt.Sprintf("%v", stage.ProcessingTime.Seconds()))
		}
		record = append(record, fmt.Sprintf("%v", stat.TotalProcessingTime.Seconds()))
		if err := writer.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	writer.Flush()
}

func(collector *Collector) SaveAsCSV(path string) {
	file, err := os.Create(path)
	shared.FailOnError(err, "Cannot create file")
	defer file.Close()

	collector.ToCSV(file)
}
