package yenta

import (
	"github.com/utubun/yenta/internal/model"
	"github.com/utubun/yenta/internal/service"
)

type Yenta interface {
	Subscribe(exchange model.Exchange, queue model.Queue, rout string, worker func(map[string]interface{}))
	Run() error
}

func New(url string) (Yenta, error) {
	s, err := service.New(url)
	if err != nil {
		return nil, err
	}

	return s, err
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

type Queue interface{}

func NewQueue(name string, durable, autoDelete, exclusive, noWait bool) (Queue, error) {
	q := model.NewQueue(name, durable, autoDelete, exclusive, noWait)
	return q, nil
}
