package yenta

import (
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	arg        amqp.Table
}

func NewQueue(name string, durable, autoDelete, exclusive, noWait bool, arg amqp.Table) Queue {
	return Queue{
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		arg,
	}
}

func (q Queue) Declare(ch *amqp.Channel) (amqp.Queue, error) {
	res, err := ch.QueueDeclare(
		q.name,
		q.durable,
		q.autoDelete,
		q.exclusive,
		q.noWait,
		q.arg,
	)
	return res, err
}

func (q Queue) Bind(ch *amqp.Channel, exchange Exchange, routing Routing) error {
	if !(exchange.durable && q.durable) {
		return errors.New("can not bind queue to exchange: both queue and exchange must be durable")
	}
	for _, key := range routing {
		err := ch.QueueBind(q.name, key, exchange.name, true, nil)
		log.Printf("Bindied the queue %s, key %s\n", q.name, key)
		if err != nil {
			return err
		}
	}
	return nil
}
