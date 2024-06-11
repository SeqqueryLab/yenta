package yenta

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func (s *Service) Publisher(exchange Exchange, keys Routing, mandatory, immideate bool) func(m interface{}) error {
	return func(message interface{}) error {
		var wg sync.WaitGroup
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
			wg.Add(1)
			go func(m interface{}, key string) {
				defer wg.Done()
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

		go func() {
			wg.Wait()
			ch.Close()
			log.Println("Publisher closed the channel")
			conn.Close()
			log.Println("Publisher closed the connection")
		}()

		return nil
	}
}
