package model

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func (s *Service) Publisher(exchange Exchange, queue Queue, keys Routing, mandatory, immideate bool) func(m interface{}) error {
	return func(message interface{}) error {
		conn := s.Connection()
		log.Println("Publisher created connection")
		ch := s.openChannel(conn)
		log.Println("Publisher opened a new channel")

		err := ch.ExchangeDeclare(
			exchange.name,
			exchange.kind,
			exchange.durable,
			exchange.autoDelete,
			exchange.internal,
			exchange.noWait,
			exchange.arg,
		)
		log.Println("Publisher declared exchange")

		if err != nil {
			log.Println("Error declaring exchange: ", err)
			*s.err <- err
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		arg, err := json.Marshal(message)
		if err != nil {
			return err
		}

		for _, key := range keys {
			go func(m interface{}, key string) {
				err = ch.PublishWithContext(
					ctx,
					exchange.name,
					key,
					mandatory,
					immideate,
					amqp091.Publishing{
						DeliveryMode: amqp091.Persistent,
						ContentType:  "application/json",
						Body:         arg,
					},
				)
				if err != nil {
					log.Printf("Publisher failed to send the message: %+v, for routing key %s", m, key)
					*s.err <- err
					return
				}
				log.Println("Publisher sent the message: ", m)

			}(message, key)
		}

		return nil
	}
}
