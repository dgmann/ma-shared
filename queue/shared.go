package queue

import (
	"log"
	"fmt"
	"github.com/streadway/amqp"
	"encoding/json"
	"github.com/dgmann/ma-shared"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

type Delivery struct {
	*amqp.Delivery
	Message *shared.Message
}

func NewDelivery(delivery *amqp.Delivery) (*Delivery) {
	return &Delivery{delivery, toMessage(delivery)}
}

func toMessage(delivery *amqp.Delivery) (*shared.Message) {
	var message shared.Message
	json.Unmarshal([]byte(delivery.Body), &message)
	return &message
}
