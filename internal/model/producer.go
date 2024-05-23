package model

type Producer struct {
	Exchange Exchange
	Queue    Queue
	Rout     string
	Worker   Worker
}
