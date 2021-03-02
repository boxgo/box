package rabbitmq_test

import (
	"github.com/boxgo/box/pkg/client/rabbitmq"
	"github.com/boxgo/box/pkg/logger"
)

func Example() {
	ch, err := rabbitmq.StdConfig("default").Build().Channel()
	if err != nil {
		logger.Panic(err)
	}

	q, err := ch.QueueDeclare("queue-name", true, true, false, false, nil)
	if err != nil {
		logger.Panic(err, q)
	}

	err = ch.Publish("", "queue-name", false, false, rabbitmq.Publishing{
		DeliveryMode: rabbitmq.Persistent,
		ContentType:  "text/plain",
		Body:         []byte("hello world"),
	})
	if err != nil {
		logger.Panic(err)
	}
}
