package service

import (
	"fmt"
	"testing"

	"github.com/utubun/yenta/internal/model"
)

func TestYenta(t *testing.T) {
	var s Service
	err := s.dial("amqp://guest:guest@localhost:5672/")

	t.Run("Test connection can be established", func(t *testing.T) {
		got := err

		if got != nil {
			t.Errorf("Error connecting the RabbitMQ: got %s, want %s", got, "nil")
		}
	})

	t.Run("Test that connection is open", func(t *testing.T) {
		got := s.conn.IsClosed()
		want := false

		if got != want {
			t.Errorf("Error testing if connection is closed: got %t, want %t", got, want)
		}
	})

	t.Run("test that channel can be opened", func(t *testing.T) {
		ch, err := s.openChannel()

		if err != nil {
			t.Errorf("Error trying to open the channel: %s", err)
		}

		defer ch.Close()
	})

	t.Run("exchange can be declared", func(t *testing.T) {
		e := model.NewExchange("test", "direct", false, false, false, false)
		err := e.Verify()
		if err != nil {
			t.Errorf("exchange verification failed: %s", err)
		}

		q := model.NewQueue("logs", false, true, false, false)

		f := func(map[string]interface{}) map[string]interface{} {
			fmt.Printf("Done")
			return nil
		}

		s.Subscribe(e, q, "user", f)
		if err != nil {
			t.Errorf("can not declare exchange: %s", err)
		}
	})

}
