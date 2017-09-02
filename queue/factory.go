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
	conn := connect(config.ToConnectionString())
	return &Factory{
		conn:conn,
	}
}

func connect(uri string) *amqp.Connection {
  for {
    conn, err := amqp.Dial(uri)

    if err == nil {
      return conn
    }

    log.Println(err)
    log.Printf("Trying to reconnect to queue at %s\n", uri)
    time.Sleep(500 * time.Millisecond)
  }
}

func(factory *Factory) Close()  {
	factory.conn.Close()
}
