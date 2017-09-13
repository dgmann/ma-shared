package stats

import (
	"encoding/csv"
	"os"
	"log"
	"fmt"
	"github.com/dgmann/ma-shared"
	"io"
	"sort"
)

type Collector struct {
	Stats Statistics
}

type Statistics []Stat

func(statistics Statistics) Len()  int {
	return len(statistics)
}

func(statistics Statistics) Less(i, j int)  bool {
	return statistics[i].Message.FrameNumber < statistics[j].Message.FrameNumber
}

func(statistics Statistics) Swap(i, j int) {
	statistics[i], statistics[j] = statistics[j], statistics[i]
}


func NewCollector() *Collector {
	return &Collector{Stats: make(Statistics, 0, 1500)}
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
	sort.Sort(collector.Stats)

	for _, stat := range collector.Stats {
		record := []string{}
		for _, stage := range stat.Stages {
			record = append(record, fmt.Sprintf("%f", stage.ProcessingTime.Seconds()))
		}
		record = append(record, fmt.Sprintf("%f", stat.TotalProcessingTime.Seconds()))
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
