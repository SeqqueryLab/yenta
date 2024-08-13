package yenta

import amqp "github.com/rabbitmq/amqp091-go"

// Exchange
type Exchange struct {
	Name       string     `json:"name"`
	Kind       string     `json:"kind"`
	Durable    bool       `json:"durable"`
	AutoDelete bool       `json:"autodelete"`
	Internal   bool       `json:"internal"`
	NoWait     bool       `json:"nowait"`
	Arg        amqp.Table `json:"arg"`
}

func (s *Service) DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool, arg amqp.Table) error {
	// open connection
	conn := s.connect()
	// open channel
	ch := s.channel(conn)
	// defer on done
	defer conn.Close()
	// declare exchange
	err := ch.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		arg,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}
	// store the config values
	s.exchange[name] = Exchange{
		Name:       name,
		Kind:       kind,
		Durable:    durable,
		AutoDelete: autoDelete,
		Internal:   internal,
		NoWait:     noWait,
		Arg:        arg,
	}
	return nil
}
