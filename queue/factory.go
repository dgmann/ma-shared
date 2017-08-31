package queue

import (
	"github.com/streadway/amqp"
	"github.com/dgmann/ma-shared"
)

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

func QueueConfigFromEnv() FactoryConfig {
	user := shared.GetEnvOrDefault("QUEUE_USER", "guest")
	password := shared.GetEnvOrDefault("QUEUE_PASSWORD", "guest")
	host := shared.GetEnvOrDefault("QUEUE_HOST", "queue")
	port := shared.GetEnvOrDefault("QUEUE_PORT", "5672")

	return FactoryConfig{
		User: user,
		Password: password,
		Host: host,
		Port: port,
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
