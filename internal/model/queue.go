package model

type Queue struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
}

func NewQueue(name string, durable, autoDelete, exclusive, noWait bool) Queue {
	return Queue{
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
	}
}

func (q *Queue) Name() string {
	return q.name
}

func (q *Queue) Durable() bool {
	return q.durable
}

func (q *Queue) AutoDelete() bool {
	return q.autoDelete
}

func (q *Queue) Exclusive() bool {
	return q.exclusive
}

func (q *Queue) NoWait() bool {
	return q.noWait
}
