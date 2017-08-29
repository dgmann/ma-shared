package worker

import (
	"github.com/dgmann/ma-shared/openalpr"
	"github.com/dgmann/ma-shared/queue"
	"fmt"
	"strings"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"log"
)

type Pool struct {
	QueueFactory *queue.Factory
	Responses chan openalpr.OpenAlprResponse
	Managers map[string]*Manager
	QueueName string
}

type Manager struct {
	NumWorkers int `json:"numWorkers"`
	Consumer *queue.Consumer
	Url string
}

func(manager *Manager) Close() {
	manager.Consumer.Close()
}

type QueueMessage struct {
	Frame string `json:"frame"`
	Index int `json:"index"`
	Origin string `json:"origin"`
	IsEmpty bool `json:"isEmpty"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func NewWorkerPool(factory *queue.Factory, queueName string) *Pool {
	pool := Pool{
		QueueFactory: factory,
		Responses: make(chan openalpr.OpenAlprResponse),
		Managers: make(map[string]*Manager),
		QueueName: queueName,
	}
	return &pool
}

func(pool *Pool) Close() {
	for _, manager := range pool.Managers {
		manager.Consumer.Close()
	}
}

func(pool *Pool) Register(Url string, NumWorker int) {
	createWorkerNum := NumWorker

	if manager, ok := pool.Managers[Url]; ok {
		if manager.NumWorkers < NumWorker {
			createWorkerNum = NumWorker - manager.NumWorkers
			fmt.Printf("Increase workers for url %s from %d to %d.\r\n", Url, manager.NumWorkers, NumWorker)
		} else if manager.NumWorkers > NumWorker {
			manager.Consumer.Close()
			consumer, err := CreateConsumer(pool.QueueFactory, NumWorker, pool.QueueName)
			if err != nil {
				fmt.Printf("Could not create channel for %s. Could not register %s\r\n", Url, Url)
				return
			}
			manager.Consumer = consumer
		} else {
			fmt.Printf("Url %s with %d workers already registered.\r\n", Url, NumWorker)
			return
		}

	} else {
		consumer, err := CreateConsumer(pool.QueueFactory, NumWorker, pool.QueueName)
		if err != nil {
			fmt.Printf("Could not create channel for %s. Could not register %s\r\n", Url, Url)
			return
		}
		pool.Managers[Url] = &Manager{
			NumWorkers: NumWorker,
			Consumer: consumer,
			Url: Url,
		}
		fmt.Printf("Registered url %s with %d workers\r\n", Url, NumWorker)
	}

	pool.Managers[Url].NumWorkers = NumWorker
	pool.Managers[Url].Start(createWorkerNum, pool.Responses)
	fmt.Printf("Url %s has %d workers\r\n", Url, pool.Managers[Url].NumWorkers)
}

func CreateConsumer(factory *queue.Factory, numWorkers int, queueName string) (*queue.Consumer, error) {
	consumer, err := factory.NewConsumer(queueName)
	consumer.Qos(numWorkers * 4, 0, false)
	return consumer, err
}

func(manager *Manager) Start(numWorkers int, responses chan<- openalpr.OpenAlprResponse) {
	msgs := manager.Consumer.Consume()
	for i := 0; i < numWorkers; i++ {
		go StartWorkerRoutine(msgs, manager.Url, responses)
	}
}

func StartWorkerRoutine(msgs <-chan *queue.Delivery, url string, responses chan<- openalpr.OpenAlprResponse) {
	fmt.Printf("Start worker for url %s\r\n", url)
	defer fmt.Printf("Shutdown worker for url %s\r\n", url)

	for msg :=  range msgs {
		var message QueueMessage
		msgBody := string(msg.Body)
		startIndex := strings.Index(msgBody, "{\"")
		msgBody = msgBody[startIndex:]

		err := json.Unmarshal([]byte(msgBody), &message)

		failOnError(err, "Failed to decode queue message")
		response, err := SendRequest(url, message.Frame)
		msg.Ack(false)
		responses <- response
	}
}

func SendRequest(Url string, Frame string) (openalpr.OpenAlprResponse, error)  {
	var openAlprResponse openalpr.OpenAlprResponse

	client := http.Client{}
	req, err := http.NewRequest("POST", Url, strings.NewReader(Frame))
	req.Header.Add("Content-Type", "text/plain")
	defer req.Body.Close()

	fmt.Printf("Send Request to %s\r\n", Url)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return openAlprResponse, err
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(bodyBytes, &openAlprResponse)
	return openAlprResponse, nil
}
