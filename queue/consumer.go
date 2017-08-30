package queue

import (
	"github.com/streadway/amqp"
)

type Consumer struct {
	channel *amqp.Channel
	queue amqp.Queue
}

func(factory *Factory) NewConsumer(queueName string) (*Consumer, error) {
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

	return &Consumer{
		channel: ch,
		queue:queue,
	}, err
}

func(consumer *Consumer) Qos(prefetchCount int, prefetchSize int, global bool) (error) {
	return consumer.channel.Qos(
		prefetchCount,
		prefetchSize,
		false,
	)
}

func(consumer *Consumer) Consume() (<- chan *Delivery) {
	del, _ := consumer.channel.Consume(
		consumer.queue.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	deliveries := make(chan *Delivery)
	go func() {
		for d := range del {
			deliveries <- NewDelivery(d)
		}
	}()

	return deliveries
}

func(consumer *Consumer) Close()  {
	consumer.channel.Close()
}
