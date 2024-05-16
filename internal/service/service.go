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

type Exchange interface {
	model.Exchange
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
		log.Fatal("Can not connect to RabbitMQ: ", err)
		return err
	}
	s.conn = conn
	return err
}

func (s *Service) openChannel() (*amqp.Channel, error) {
	return s.conn.Channel()
}

func (s *Service) Subscribe(exchange model.Exchange, queue model.Queue, rout string, worker func(map[string]interface{})) {
	go subscribe(s, *s.err, exchange, queue, rout, worker)()
	log.Printf("subscribed to the exchange: %s, que: %s, binding key: %s\n", exchange.Name(), queue.Name(), rout)
}

func (s *Service) Run() error {
	log.Printf("Running the service...")
	return s.err
}
