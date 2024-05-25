package yenta

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/utubun/yenta/internal/model"
)

type Service interface {
	Run() model.YentaError
	Worker(exchange model.Exchange, queue model.Queue, key model.Routing, fun func(arg interface{}) interface{})
}

func NewService(url string) Service {
	return model.NewService(url)
}

type Exchange interface{}

func NewExchange(name, kind string, durable, autoDelete, internal, noWait bool, arg amqp.Table) Exchange {
	exchange := model.NewExchange(name, kind, durable, autoDelete, internal, noWait, arg)
	return exchange
}

type Queue interface {
	Declare() (amqp.Queue, error)
	Bind(ch *amqp.Channel, exchange *Exchange, routing *Routing) error
}

func NewQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) Queue {
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

/*
func publish(s *Service, exchange Exchange, queue Queue, rout string, worker Worker, in chan interface{}) {
	conn := s.Connection()
	ch := s.openChannel(conn)

	err := ch.ExchangeDeclare(
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
				rout,
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
*/
