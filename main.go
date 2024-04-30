package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Envoy struct{}

func (env *Envoy) Receive(path string) error {
	conn, err := amqp.Dial(path)
	if err != nil {
		log.Fatalf("Connection to RabbitMQ failed: %s", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel: %s", err)
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
		return err
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
		log.Fatalf("Failed to consume message: %s", err)
		return err
	}

	var forever chan struct{}

	go func() {
		for m := range msg {
			log.Printf("Received message: %s", m.Body)
		}
	}()
	log.Println("Waiting for the messages from sender")
	<-forever
	return err
}

func (env *Envoy) Send(path string) {
	con, err := amqp.Dial(path)
	if err != nil {
		log.Fatal("Connection error")
	}
	ch, err := con.Channel()
	if err != nil {
		log.Printf("Can not open the channel: %s", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Can not declare the queue: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var forever chan struct{}

	go func() {
		for {
			body := "Hi!"
			err = ch.PublishWithContext(
				ctx,
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				},
			)
			if err != nil {
				log.Printf("Failed to publish message: %s", err)
				return
			}
			log.Printf("[+] Sent message: %s", body)
			time.Sleep(5 * time.Second)
		}
	}()
	<-forever
}

func main() {
	path := "amqp://guest:guest@localhost:5672/"
	env := &Envoy{}

	go env.Send(path)
	go env.Receive(path)
	time.Sleep(30 * time.Second)
}
