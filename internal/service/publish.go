package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/utubun/yenta/internal/model"
)

func Send(s *Service, err model.YentaError, exchange model.Exchange, queue model.Queue, rout string, worker func(map[string]interface{})) func() {
	return func() {
		c, e := s.openChannel()
		if err != nil {
			err <- e
		}
		defer c.Close()

		q, e := c.QueueDeclare(
			queue.Name(),
			queue.Durable(),
			queue.AutoDelete(),
			queue.Exclusive(),
			queue.NoWait(),
			nil,
		)
		if e != nil {
			err <- e
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		go func() {
			for {
				body := make(map[string]interface{})
				body["message"] = "Hi!"
				body["time"] = time.Now()
				args, _ := json.Marshal(body)
				e = c.PublishWithContext(
					ctx,
					"",
					q.Name,
					false,
					false,
					amqp.Publishing{
						ContentType: "application/json",
						Body:        args,
					},
				)
				if e != nil {
					err <- e
				}
				log.Println("Message sent!")
				time.Sleep(5000 * time.Millisecond)

			}
		}()
	}
}
