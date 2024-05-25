package model

import (
	"log"
)

func Worker(exchange Exchange, queue Queue, key Routing, fun func(arg interface{}) interface{}) func(s *Service) error {
	return func(s *Service) error {
		conn := s.Connection()
		log.Print("Worker created connection")

		ch := s.openChannel(conn)
		log.Println("Worker opened the channel")

		err := exchange.Declare(ch)
		if err != nil {
			return err
		}
		log.Println("worker created exchange")

		q, err := queue.Declare(ch)
		log.Println(err)
		if err != nil {
			return err
		}
		log.Printf("worker created queue %+v", q)

		err = queue.Bind(ch, exchange, key)
		if err != nil {
			return err
		}
		log.Println("Worker binded the queue to exchange")

		err = ch.Qos(
			1,
			0,
			false,
		)
		if err != nil {
			return err
		}

		msg, err := ch.Consume(
			queue.name,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}
		log.Println("Consuming the messages")

		for m := range msg {
			ack := m.Ack
			go func() {
				fun(m)
				ack(false)
			}()
		}

		return nil
	}
}
