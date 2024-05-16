package model

import (
	"errors"
	"fmt"
	"regexp"
)

type Exchange struct {
	name       string
	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
}

func NewExchange(name, kind string, durable, autoDelete, internal, noWait bool) Exchange {
	return Exchange{
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
	}
}
func (e *Exchange) Verify() error {
	if e.name == "" {
		return errors.New("can not declare exchange with empty name")
	}

	if match, _ := regexp.MatchString("^amq\\.?.*", e.name); match {
		return fmt.Errorf("can not declare exchange. %s is reserved name", e.name)
	}

	if ok := e.kind == "direct" || e.kind == "fanout" || e.kind == "topic" || e.kind == "headers"; !ok {
		return errors.New("Exchange kind must be 'direct', 'fanout', 'topic' or 'headers'")
	}

	return nil
}

func (e *Exchange) Name() string {
	return e.name
}

func (e *Exchange) Kind() string {
	return e.kind
}

func (e *Exchange) Durable() bool {
	return e.durable
}

func (e *Exchange) AutoDelete() bool {
	return e.autoDelete
}

func (e *Exchange) Internal() bool {
	return e.internal
}

func (e *Exchange) NoWait() bool {
	return e.noWait
}
