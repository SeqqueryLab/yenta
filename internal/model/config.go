package model

type WorkItem struct {
	Consumer Consumer
	Producer Producer
}

type Config []WorkItem
