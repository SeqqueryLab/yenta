package yenta

import amqp "github.com/rabbitmq/amqp091-go"

type Exchange struct {
	name       string
	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	arg        amqp.Table
}

func NewExchange(name, kind string, durable, autoDelete, internal, noWait bool, arg amqp.Table) Exchange {
	return Exchange{
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		arg,
	}
}

func (e *Exchange) Declare(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		e.name,
		e.kind,
		e.durable,
		e.autoDelete,
		e.internal,
		e.noWait,
		e.arg,
	)
	return err
}
