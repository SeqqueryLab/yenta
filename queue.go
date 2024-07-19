package yenta

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Queue
type Queue struct {
	Exchange   string
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Arg        amqp.Table
	Keys       []string
}

// Declare Queue
func (s *Service) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool, arg amqp.Table) error {
	// open connection
	conn := s.connect()
	// open channel
	ch := s.channel(conn)
	// defer on done
	defer conn.Close()
	// Declare the queue
	_, err := ch.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		arg,
	)
	if err != nil {
		return fmt.Errorf("can not declare the queue %s", err)
	}
	s.queue[name] = Queue{
		Name:       name,
		Durable:    durable,
		AutoDelete: autoDelete,
		Exclusive:  exclusive,
		NoWait:     noWait,
		Arg:        arg,
	}
	return nil
}

// Bind Queue
func (s *Service) BindQueue(name, exchange string, keys []string, noWait bool, arg amqp.Table) error {
	// open connection
	conn := s.connect()
	//open channel
	ch := s.channel(conn)
	// defer on done
	defer conn.Close()
	// iterate over keys, and bind the queue to exchange
	for _, key := range keys {
		err := ch.QueueBind(name, key, exchange, noWait, arg)
		if err != nil {
			return err
		}
	}
	return nil
}
