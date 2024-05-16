package model

import "testing"

func TestExchange(t *testing.T) {
	t.Run("Exchange can be configured and verified", func(t *testing.T) {
		e := NewExchange("log", "topic", false, false, false, true)

		if err := e.Verify(); err != nil {
			t.Errorf("exchange verification failed: %s", err)
		}
	})

	t.Run("Exchange can not be configured with empty name", func(t *testing.T) {
		e := NewExchange("", "topic", false, false, false, true)
		if err := e.Verify(); err == nil {
			t.Error("exchange should not be created with empty name")
		}
	})

	t.Run("Exchange can not be configured with name starting from 'amq.'", func(t *testing.T) {
		e := NewExchange("amq.", "topic", false, false, false, true)
		if err := e.Verify(); err == nil {
			t.Error("exchange should not be configured with name starting from 'amq.'")
		}
	})
}
