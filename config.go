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
	Exchange string     `json:"exchange"`
	Queue    string     `json:"queue"`
	Keys     []string   `json:"keys"`
	NoWait   bool       `json:"nowait"`
	Arg      amqp.Table `json:"arg"`
}

// ConfigJSON
func configFromJSON(b []byte) Config {
	var res Config
	json.Unmarshal(b, &res)
	return res
}
