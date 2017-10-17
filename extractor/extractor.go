package extractor

import (
	"github.com/dgmann/ma-shared/openalpr"
	"github.com/dgmann/ma-shared"
	"sync"
)

type ExtractorFactory struct {
	country string
	config string
	runtimeDataDir string
	alpr []*openalpr.Alpr
}

func NewExtractorFactory(country, config, runtimeDataDir string) ExtractorFactory {
	return ExtractorFactory{
		country: country,
		config: config,
		runtimeDataDir: runtimeDataDir,
	}
}

func(factory *ExtractorFactory) Initialize(numExtractors int) {
	var instances []*openalpr.Alpr
	for i:=0; i < numExtractors; i++ {
		instances = append(instances, openalpr.NewAlpr(factory.country, factory.country, factory.runtimeDataDir))
	}
	factory.alpr = instances
}

func(factory *ExtractorFactory) StartExtractors(input <-chan shared.Message) chan shared.Message {
	var wg sync.WaitGroup
	output := make(chan shared.Message, 10000)

	for i:=0; i < len(factory.alpr); i++ {
		wg.Add(1)
		go factory.startExtractor(factory.alpr[i], input, output, &wg)
	}
	go func() {
		wg.Wait()
		close(output)
		for _, oalpr := range factory.alpr {
			oalpr.Unload()
		}
	}()
	return output
}

func(factory *ExtractorFactory) startExtractor(oalpr *openalpr.Alpr, input <-chan shared.Message, output chan<-shared.Message, wg *sync.WaitGroup) {
	for msg := range input {
		msg.EnterStage("Extractor")
		result, _ := oalpr.RecognizeByBlob(msg.Image)
		msg.Result.OpenALPR = result
		msg.LeaveStage("Extractor")
		output <- msg
	}

	wg.Done()
}
