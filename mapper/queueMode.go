package mapper

import (
	"github.com/dgmann/ma-shared/queue"
	"runtime"
	"github.com/dgmann/ma-shared"
	"sync"
	"github.com/jeffail/tunny"
	"fmt"
)

type QueueMode struct {
	Mode
	FactoryConfig queue.FactoryConfig
	OutputQueueName string
	PrefetchCount int
	WorkerCount int
}

func NewQueueMode(stageName string, pool *tunny.WorkPool, factoryConfig queue.FactoryConfig, outputQueueName string) *QueueMode {
	return &QueueMode{
		Mode{stageName, pool},
		factoryConfig,
		outputQueueName,
		runtime.NumCPU(),
		runtime.NumCPU(),
	}
}

func(mode *QueueMode) Listen(inputQueue string, setResult func(message *shared.Message, result interface{})) {
	println("Listening on " + mode.FactoryConfig.ToConnectionString())
	factory := queue.NewFactory(mode.FactoryConfig)
	consumer, _ := factory.NewConsumer(inputQueue)
	consumer.Qos(mode.PrefetchCount, 0, false)
	deliveries := consumer.Consume()
	var wg sync.WaitGroup

	fmt.Printf("Starting %v workers\r\n", mode.WorkerCount)
	for i := 0; i < mode.WorkerCount; i++ {
		go func() {
			wg.Add(1)
			println("Worker started")
			producer := factory.NewProducer(mode.OutputQueueName)
			for delivery := range deliveries {
				message := delivery.Message
				message.EnterStage(mode.StageName)
				result, _ := mode.pool.SendWork(message)
				setResult(message, result)
				message.LeaveStage(mode.StageName)
				producer.SendAsJSON(message)
				delivery.Ack(false)
			}
			producer.Close()
			wg.Done()
			println("Worker stopped")
		}()
	}
	wg.Wait()
}
