package yenta

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/utubun/yenta/internal/model"
)

type Service struct {
	conn *amqp.Connection
	err  *model.YentaError
}

func New(url string) (*Service, error) {
	var s Service
	err := s.dial(url)
	if err != nil {
		return nil, err
	}

	e := make(model.YentaError)
	s.err = &e

	return &s, err
}

func (s *Service) dial(url string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	s.conn = conn
	return err
}

func (s *Service) openChannel() (*amqp.Channel, error) {
	return s.conn.Channel()
}

type Exchange interface {
	Verify() error
	Name() string
	Kind() string
	Durable() bool
	AutoDelete() bool
	Internal() bool
	NoWait() bool
}

func NewExchange(name, kind string, durable, autoDelete, internal, noWait bool) (*model.Exchange, error) {
	e := model.NewExchange(name, kind, durable, autoDelete, internal, noWait)
	if err := e.Verify(); err != nil {
		return nil, err
	}
	return &e, nil
}

type Queue interface {
	Name() string
	Durable() bool
	AutoDelete() bool
	Exclusive() bool
	NoWait() bool
}

func NewQueue(name string, durable, autoDelete, exclusive, noWait bool) (*model.Queue, error) {
	q := model.NewQueue(name, durable, autoDelete, exclusive, noWait)
	return &q, nil
}

type Producer struct {
	Exchange Exchange
	Queue    Queue
	Rout     string
	Worker   Worker
}

type Consumer struct {
	Exchange Exchange
	Queue    Queue
	Rout     string
	Worker   Worker
}

type Worker func(map[string]interface{}) map[string]interface{}

func (w Worker) Do(ch chan interface{}, args map[string]interface{}) {
	go func() {
		ch <- w(args)
	}()
}

type Config []struct {
	Consumer Consumer
	Producer Producer
}

func (s *Service) Add(config Config) {
	log.Println("Routing the service")
	for _, w := range config {
		out := subscribe(s, w.Consumer.Exchange, w.Consumer.Queue, w.Consumer.Rout, w.Consumer.Worker)
		log.Println("Consumer is ready")
		publish(s, w.Producer.Exchange, w.Producer.Queue, w.Producer.Rout, w.Producer.Worker, out)
		log.Println("Producer is ready")
	}
	log.Println("Routing completed")
}

func (s *Service) Run() error {
	log.Printf("Running the service...")
	return s.err
}

func subscribe(s *Service, exchange Exchange, queue Queue, rout string, worker Worker) chan interface{} {
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

func publish(s *Service, exchange Exchange, queue Queue, rout string, worker Worker, in chan interface{}) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	out := make(chan interface{})

	go func() {
		for arg := range in {
			worker.Do(out, arg.(map[string]interface{}))
			log.Printf("Producer received the message: %+v", arg)
		}
	}()

	go func() {
		for m := range out {
			args, _ := json.Marshal(m)
			err = ch.PublishWithContext(
				ctx,
				exchange.Name(),
				q.Name,
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "application/json",
					Body:         args,
				},
			)
			if err != nil {
				log.Printf("Error publishing the message: %s\n", err)
			}
			log.Println("Message emitted")
		}
	}()
}
