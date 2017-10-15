package lookup

import (
	"github.com/dgmann/ma-shared"
	"sync"
)

type LookupFactory struct {
	client StorageClient
	numWorker int
}

func NewLookupFactory(client StorageClient, numWorker int) *LookupFactory {
	return &LookupFactory{client:client, numWorker:numWorker}
}

func(factory *LookupFactory) StartLookupServices(input <-chan shared.Message) chan shared.Message {
	var wg sync.WaitGroup
	output := make(chan shared.Message, 10000)
	for i:=0; i < factory.numWorker; i++ {
		wg.Add(1)
		go factory.lookup(input, output, &wg)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func(factory *LookupFactory) lookup(input <-chan shared.Message, output chan<-shared.Message, wg *sync.WaitGroup) {
	for msg := range input {
		msg.EnterStage("Lookup")
		for _, result := range msg.Result.OpenALPR.Results {
			exists := factory.client.Exists(result.Plate)
			if exists {
				msg.Result.WantedNumbers = append(msg.Result.WantedNumbers, result.Plate)
				println(result.Plate)
			}
		}
		msg.LeaveStage("Lookup")
		output <- msg
	}
	wg.Done()
}
