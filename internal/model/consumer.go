package model

type Consumer struct {
	Exchange Exchange
	Queue    Queue
	Rout     string
	Worker   Worker
}
