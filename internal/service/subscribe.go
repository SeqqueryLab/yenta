package service

import (
	"encoding/json"
	"log"

	"github.com/utubun/yenta/internal/model"
)

func subscribe(s *Service, exchange model.Exchange, queue model.Queue, rout string, worker model.Worker) chan interface{} {
	ch, err := s.openChannel()
	if err != nil {
		*s.err <- err
	}
	log.Println("Channel created...")

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
	log.Printf("Exchange %s declared...\n", exchange.Name())

	q, err := ch.QueueDeclare(
		queue.Name(),
		queue.Durable(),
		queue.AutoDelete(),
		queue.Exclusive(),
		queue.NoWait(),
		nil,
	)
	if err != nil {
		log.Printf("Error decalring the queue: %s\n", err)
		*s.err <- err
	}
	log.Printf("Queue %s declared\n", queue.Name())

	if rout == "" {
		rout = "#"
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
	log.Printf("Queue %s binded to exchange %s\n", queue.Name(), exchange.Name())

	err = ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		*s.err <- err
	}

	msg, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		*s.err <- err
	}
	log.Printf("Consuming the messages for the queue %s\n", queue.Name())

	out := make(chan interface{})

	go func() {
		for m := range msg {
			args := make(map[string]interface{})
			json.Unmarshal(m.Body, &args)
			worker.Do(out, args)
		}
		defer log.Println("Stop sending")
		defer ch.Close()
	}()

	return out

}
