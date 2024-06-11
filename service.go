package yenta

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Service struct {
	url string
	err *YentaError
}

func NewService(url string) *Service {
	var s Service

	s.url = url
	err := make(YentaError)
	s.err = &err

	return &s
}

func (s *Service) Connection() *amqp.Connection {
	conn, err := amqp.Dial(s.url)
	if err != nil {
		*s.err <- err
	}

	return conn
}

func (s *Service) openChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		*s.err <- err
	}
	return ch
}

func (s *Service) Worker(exchange Exchange, queue Queue, key Routing, fun func(arg interface{}) interface{}) {
	w := Worker(exchange, queue, key, fun)
	log.Printf("Starting the worker")
	go w(s)
}

func (s *Service) Run() YentaError {
	return *s.err
}
