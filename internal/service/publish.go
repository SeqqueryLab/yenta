package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/utubun/yenta/internal/model"
)

func publish(s *Service, exchange model.Exchange, queue model.Queue, rout string, worker model.Worker, in chan interface{}) {
	ch, err := s.openChannel()
	if err != nil {
		*s.err <- err
	}

	err = ch.ExchangeDeclare(
		exchange.Name(),
		exchange.Kind(),
		exchange.Durable(),
		exchange.AutoDelete(),
		exchange.Internal(),
		exchange.NoWait(),
		nil,
	)
	if err != nil {
		*s.err <- err
	}

	q, err := ch.QueueDeclare(
		queue.Name(),
		queue.Durable(),
		queue.AutoDelete(),
		queue.Exclusive(),
		queue.NoWait(),
		nil,
	)
	if err != nil {
		*s.err <- err
	}

	err = ch.QueueBind(
		q.Name,
		rout,
		exchange.Name(),
		queue.NoWait(),
		nil,
	)
	if err != nil {
		*s.err <- err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var out chan interface{}

	go func() {
		for arg := range in {
			worker.Do(out, arg.(map[string]interface{}))
			log.Printf("Producer received the message: %+v", arg)
		}
	}()

	go func() {
		for m := range out {
			args, _ := json.Marshal(m)
			log.Printf("Produced the message: %s\n", m)
			err = ch.PublishWithContext(
				ctx,
				exchange.Name(),
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        args,
				},
			)
			if err != nil {
				log.Printf("Error publishing the message: %s\n", err)
			}
		}
	}()
}
