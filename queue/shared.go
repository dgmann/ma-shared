package queue

import (
	"log"
	"fmt"
	"time"
	"github.com/streadway/amqp"
	"encoding/json"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

type Delivery struct {
	*amqp.Delivery
	Message *Message
}

func NewDelivery(delivery *amqp.Delivery) (*Delivery) {
	return &Delivery{delivery, toMessage(delivery)}
}

func(message *Message) EnterStage(stageName string) {
	message.Stages[stageName] = Stage{
		EnteredAt: time.Now(),
	}
}

func(message *Message) LeaveStage(stageName string) {
	message.Stages[stageName] = Stage{
		EnteredAt: message.Stages[stageName].EnteredAt,
		LeftAt: time.Now(),
	}
}

func toMessage(delivery *amqp.Delivery) (*Message) {
	var message Message
	json.Unmarshal([]byte(delivery.Body), &message)
	return &message
}

type Message struct {
	Image []byte `json:"image"`
	FrameNumber int `json:"frameNumber"`
	CreatedAt time.Time `json:"createdAt"`
	Stages map[string]Stage `json:"timeline"`
}

type Stage struct {
	EnteredAt time.Time `json:"enteredAt"`
	LeftAt time.Time `json:"leftAt"`
}

func NewMessage(image []byte, frameNumer int, createdAt time.Time) (*Message) {
	return & Message{
		Image: image,
		FrameNumber: frameNumer,
		CreatedAt: createdAt,
		Stages: make(map[string]Stage),
	}
}
