package stats

import (
	"github.com/dgmann/ma-shared"
	"time"
	"sort"
	"text/tabwriter"
	"os"
	"fmt"
	"strconv"
)

func NewStat(message *shared.Message) (Stat) {
	stages := calculateProcessingTimes(message)
	return Stat{
		message,
		time.Now(),
		stages[0].EnteredAt.Sub(stages[len(stages)-1].LeftAt),
		stages,
	}
}

func(stat *Stat) Print()  {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, '-', tabwriter.Debug)
	values := ""
	for _, stage := range stat.Stages {
		values = values + stage.Name + ": " + strconv.Itoa(int(stage.ProcessingTime.Nanoseconds() / 1000000)) + "ms\t"
	}
	values = values + "Total: " + strconv.Itoa(int(stat.TotalProcessingTime.Nanoseconds() / 1000000)) + "ms\t"
	fmt.Fprintln(w, values)
	w.Flush()
}

type Stat struct {
	*shared.Message
	ReceivedAt time.Time
	TotalProcessingTime time.Duration
	Stages Stages
}

type Stage struct {
	shared.Stage
	Name string
	ProcessingTime time.Duration
}

type Stages []Stage

func(stages Stages) Len()  int {
	return len(stages)
}

func(stages Stages) Less(i, j int)  bool {
	return stages[i].EnteredAt.Before(stages[j].EnteredAt)
}

func(stages Stages) Swap(i, j int) {
	stages[i], stages[j] = stages[j], stages[i]
}

func calculateProcessingTimes(message *shared.Message) (Stages) {
	stages := Stages{}
	for k, st := range message.Stages {
		stages = append(stages, Stage{Stage: st, Name: k, ProcessingTime: st.LeftAt.Sub(st.EnteredAt)})
	}
	sort.Sort(stages)

	times := Stages{}
	for i := 0; i < len(stages); i++ {
		current := stages[i]
		current.ProcessingTime = current.LeftAt.Sub(current.EnteredAt)
		if i == len(stages) - 1 {
			times = append(times, current)
			break
		}
		intermediate := Stage{
			Name: current.Name + "<->" + stages[i+1].Name,
			Stage: shared.Stage{
				EnteredAt: current.LeftAt,
				LeftAt: stages[i+1].EnteredAt,
			},
		}
		intermediate.ProcessingTime = intermediate.LeftAt.Sub(intermediate.EnteredAt)
		times = append(times, current, intermediate)
	}
	return times
}
