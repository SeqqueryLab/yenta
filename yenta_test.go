package yenta

import (
	"fmt"
	"testing"
)

func TestYenta(t *testing.T) {
	s, err := New("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Errorf("can not initialize Yenta: %s", err)
	}
	t.Run("Yenta is created and and channel can be opened", func(t *testing.T) {
		// placeholder
		fmt.Printf("service connection: %s", s)
	})

	t.Run("Exchange can be configured and verified", func(t *testing.T) {
		_, err := NewExchange("log", "topic", false, false, false, true)
		if err != nil {
			t.Errorf("exchange verification failed: %s", err)
		}
	})

	t.Run("Exchange can not be configured with empty name", func(t *testing.T) {
		_, err := NewExchange("", "topic", false, false, false, true)
		if err == nil {
			t.Error("exchange should not be created with empty name")
		}
	})

	t.Run("Exchange can not be configured with name starting from 'amq.'", func(t *testing.T) {
		_, err := NewExchange("amq.", "topic", false, false, false, true)
		if err == nil {
			t.Error("exchange should not be configured with name starting from 'amq.'")
		}
	})
}
