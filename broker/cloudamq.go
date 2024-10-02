package broker

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CloudAMQ struct {
	URL       string
	QueueName string
}

func NewCloudAMQ(url string, queue string) *CloudAMQ {
	return &CloudAMQ{
		URL:       url,
		QueueName: queue,
	}
}

func (c CloudAMQ) Send(message string) error {
	connection, err := amqp.Dial(c.URL)
	if err != nil {
		return err
	}
	defer connection.Close()

	ch, err := connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(c.QueueName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		Timestamp:   time.Now(),
		ContentType: "text/plain",
		Body:        []byte(message),
	}

	returns := ch.NotifyReturn(make(chan amqp.Return))
	go func() {
		for r := range returns {
			fmt.Printf("Warning: Message returned from RabbitMQ: %s", r.ReplyText)
		}
	}()

	mandatory, immediate := true, false
	err = ch.Publish("", queue.Name, mandatory, immediate, msg)
	if err != nil {
		return err
	}

	return nil
}
