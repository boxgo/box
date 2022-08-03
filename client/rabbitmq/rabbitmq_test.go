package rabbitmq_test

import (
	rabbitmq2 "github.com/boxgo/box/v2/client/rabbitmq"
	"github.com/boxgo/box/v2/logger"
)

func Example() {
	ch, err := rabbitmq2.StdConfig("default").Build().Channel()
	if err != nil {
		logger.Panic(err)
	}

	q, err := ch.QueueDeclare("queue-name", true, true, false, false, nil)
	if err != nil {
		logger.Panic(err, q)
	}

	err = ch.Publish("", "queue-name", false, false, rabbitmq2.Publishing{
		DeliveryMode: rabbitmq2.Persistent,
		ContentType:  "text/plain",
		Body:         []byte("hello world"),
	})
	if err != nil {
		logger.Panic(err)
	}
}
