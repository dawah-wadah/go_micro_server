package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         //  is this exchange durable?
		false,        //  is this exchange auto-deleted?
		false,        //  is this exchange exclusive?
		false,        //  is this exchange no-local? no wait?
		nil,          //  arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {

	return ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // is this exchange
		false, // no-wait
		nil,   // arguments
	)
}
