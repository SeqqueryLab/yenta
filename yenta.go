package yenta

import (
	"github.com/SeqqueryLab/yenta/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Service interface {
	Run() model.YentaError
	Worker(exchange model.Exchange, queue model.Queue, key model.Routing, fun func(arg interface{}) interface{})
	Publisher(exchange model.Exchange, queue model.Queue, keys model.Routing, mandatory, immideate bool) func(message interface{}) error
}

func NewService(url string) Service {
	return model.NewService(url)
}

type Exchange interface{}

func NewExchange(name, kind string, durable, autoDelete, internal, noWait bool, arg amqp.Table) model.Exchange {
	exchange := model.NewExchange(name, kind, durable, autoDelete, internal, noWait, arg)
	return exchange
}

type Queue interface {
	Declare() (amqp.Queue, error)
	Bind(ch *amqp.Channel, exchange *Exchange, routing *Routing) error
}

func NewQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) model.Queue {
	q := model.NewQueue(name, durable, autoDelete, exclusive, noWait, args)
	return q
}

type Routing interface{}

func NewRouting(keys []string) Routing {
	return &model.Routing{}
}

type Message struct {
	amqp.Delivery
}
