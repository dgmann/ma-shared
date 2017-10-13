package stats

import "github.com/dgmann/ma-shared"

type StatCreator struct {

}

func NewStatCreator() StatCreator {
	return StatCreator{}
}

func(factory *StatCreator) Start(input <-chan shared.Message) chan Stat {
	output := make(chan Stat, 10000)
	go func() {
		for msg := range input {
			stat := NewStat(&msg)
			output <- stat
		}
		close(output)
	}()
	return output
}