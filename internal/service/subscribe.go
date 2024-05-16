package service

import (
	"encoding/json"
	"log"

	"github.com/utubun/yenta/internal/model"
)

func subscribe(s *Service, err model.YentaError, exchange model.Exchange, queue model.Queue, rout string, worker func(map[string]interface{})) func() {
	return func() {
		c, e := s.openChannel()
		if e != nil {
			err <- e
		}

		e = c.ExchangeDeclare(
			exchange.Name(),
			exchange.Kind(),
			exchange.Durable(),
			exchange.AutoDelete(),
			exchange.Internal(),
			exchange.NoWait(),
			nil,
		)
		if e != nil {
			err <- e
		}

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

		if rout == "" {
			rout = "#"
		}
		e = c.QueueBind(
			q.Name,
			rout,
			exchange.Name(),
			queue.NoWait(),
			nil,
		)
		if e != nil {
			err <- e
		}

		e = c.Qos(
			1,
			0,
			false,
		)
		if e != nil {
			err <- e
		}

		msg, e := c.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if e != nil {
			err <- e
		}

		go func() {
			for m := range msg {
				args := make(map[string]interface{})
				json.Unmarshal(m.Body, &args)
				go func() {
					worker(args)
				}()
			}
			defer log.Println("Stop sending")
			defer c.Close()
		}()
	}
}
