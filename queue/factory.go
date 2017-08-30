package queue

import "github.com/streadway/amqp"

type Factory struct {
	conn *amqp.Connection
}

type FactoryConfig struct {
	Host string
	User string
	Password string
	Port string
}

func(config FactoryConfig) ToConnectionString() string {
	return "amqp://" + config.User + ":" + config.Password + "@" + config.Host + ":" + config.Port
}

func DefaultConfig() FactoryConfig {
	return FactoryConfig{
		Host: "queue",
		User: "guest",
		Password: "guest",
		Port: "5672",
	}
}

func NewFactory(config FactoryConfig) (*Factory) {
	conn, err := amqp.Dial(config.ToConnectionString())
	failOnError(err, "Failed to connect to RabbitMQ")
	return &Factory{
		conn:conn,
	}
}

func(factory *Factory) Close()  {
	factory.conn.Close()
}
