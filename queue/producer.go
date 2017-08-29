package queue

import (
	"github.com/streadway/amqp"
	"encoding/json"
)

type Producer struct {
	channel *amqp.Channel
	queue amqp.Queue
}

func(factory *Factory) NewProducer(queueName string) (*Producer) {
	ch, err := factory.conn.Channel()
	failOnError(err, "Failed to open a channel")
	queue, err := ch.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare Queue")

	return &Producer{
		channel: ch,
		queue:queue,
	}
}

func(queue *Producer) Close()  {
	queue.channel.Close()
}

func(queue *Producer) SendAsJSON(data interface{}) {
	bytes, err := json.Marshal(data)
	failOnError(err, "Failed to convert to JSON")
	queue.Send("application/json", bytes)
}

func(queue *Producer) Send(contentType string, data []byte) {
	queue.channel.Publish(
		"",     // exchange
		queue.queue.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: contentType,
			Body:        data,
		})
}
