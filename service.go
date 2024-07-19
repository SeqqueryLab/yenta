package yenta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Yenta
type Service struct {
	URL      string
	exchange map[string]Exchange
	queue    map[string]Queue
	err      Error
}

// New
func New(url string) (*Service, error) {
	err := testURL(url)
	srv := &Service{
		URL:      url,
		exchange: make(map[string]Exchange),
		queue:    make(map[string]Queue),
		err:      make(Error),
	}
	return srv, err
}

// Configure
func (s *Service) Configure(path string) error {
	// define config variable
	var config *Config
	// read config file
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// read the data
	json.Unmarshal(b, &config)
	// check that exchange is not empty
	if len(config.Exchange) == 0 {
		return errors.New("exchange must be an array of length > 0")
	}
	// declare exchange
	for _, e := range config.Exchange {
		err := s.DeclareExchange(e.Name, e.Kind, e.Durable, e.AutoDelete, e.Internal, e.NoWait, e.Arg)
		if err != nil {
			return fmt.Errorf("error declaring the exchange %s: %s", e.Name, err)
		}
	}
	// declare queues
	if len(config.Queue) != 0 {
		for _, q := range config.Queue {
			err := s.DeclareQueue(q.Name, q.Durable, q.AutoDelete, q.Exclusive, q.NoWait, q.Arg)
			if err != nil {
				return fmt.Errorf("error declaring the queue %s: %s", q.Name, err)
			}
			// bind the queue to exhange if defined
			if len(config.Binding) > 0 {
				for _, binding := range config.Binding {
					if q.Name == binding.Queue {
						s.BindQueue(binding.Queue, binding.Exchange, binding.Keys, binding.NoWait, binding.Arg)
					}
				}
			}
		}
	} else {
		fmt.Println("warning, no queues to declare")
	}
	return nil
}

// connect
func (s *Service) connect() *amqp.Connection {
	conn, err := amqp.Dial(s.URL)
	if err != nil {
		s.err <- err
	}
	return conn
}

// channel
func (s *Service) channel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		s.err <- err
	}
	return ch
}

func (s *Service) Run() error {
	return s.err
}

func (s *Service) Post(exchange string, keys []string, mandatory, immediate bool, contentType string, body []byte) error {
	// open the connection
	conn := s.connect()
	// create channel
	ch := s.channel(conn)
	// close on done
	defer ch.Close()
	defer conn.Close()
	// create context and cancell
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, key := range keys {
		log.Printf("Sending the message %s with the routing key %s to the exchange %s\n", string(body), key, exchange)
		err := ch.PublishWithContext(
			ctx,
			exchange,
			key,
			mandatory,
			immediate,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  contentType,
				Body:         body,
			},
		)
		if err != nil {
			log.Printf("Error sending the message: %s\n", err)
			return err
		}
	}
	return nil
}

func (s *Service) Get(consumer, queue string, autoack, exclusive, nolocal, nowait bool, arg amqp.Table, fun func(m []byte) any) error {
	// open the connection
	conn := s.connect()
	// create channel
	ch := s.channel(conn)
	// configure queue prefetch
	err := ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return err
	}
	// consume the messages
	msg, err := ch.Consume(
		queue,
		consumer,
		autoack,
		exclusive,
		nolocal,
		nowait,
		arg,
	)
	if err != nil {
		return fmt.Errorf("failed to consume from queue %s: %s", queue, err)
	}
	log.Println("start to receive the messages")

	for m := range msg {
		body := m.Body
		go fun(body)
	}

	log.Println("exit")
	return nil
}
