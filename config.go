package yenta

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config
type Config struct {
	Exchange []Exchange `json:"exchange"`
	Queue    []Queue    `json:"queue"`
	Binding  []Binding  `json:"binding"`
}

// Binding
type Binding struct {
	Exchange string
	Queue    string
	Keys     []string
	NoWait   bool
	Arg      amqp.Table
}

// ConfigJSON
func configFromJSON(b []byte) Config {
	var res Config
	json.Unmarshal(b, &res)
	return res
}
