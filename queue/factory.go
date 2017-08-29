package queue

import "github.com/streadway/amqp"

type Factory struct {
	conn *amqp.Connection
}

func NewFactory(url string) (*Factory) {
	conn, err := amqp.Dial("amqp://guest:guest@" + url + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	return &Factory{
		conn:conn,
	}
}

func(factory *Factory) Close()  {
	factory.conn.Close()
}
