package service

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/utubun/yenta/internal/model"
)

type Service struct {
	conn *amqp.Connection
	err  *model.YentaError
}

type Consumer struct {
	Exchange model.Exchange
	Queue    model.Queue
	Rout     string
	Worker   model.Worker
}

type Producer struct {
	Exchange model.Exchange
	Queue    model.Queue
	Rout     string
	Worker   model.Worker
}

type WorkItem struct {
	Consumer Consumer
	Producer Producer
}

type Config []WorkItem

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
	log.Println("Dialing the RabbitMQ service...")
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Can not connect to RabbitMQ: ", err)
		return err
	}
	s.conn = conn
	return err
}

func (s *Service) openChannel() (*amqp.Channel, error) {
	log.Println("Opening the channel...")
	return s.conn.Channel()
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
