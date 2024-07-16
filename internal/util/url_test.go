package util

import "testing"

func TestURL(t *testing.T) {
	t.Run("test that correct url is accepted", func(t *testing.T) {
		url := "amqp://guest:guest@amqp.dev.myorg.uk:4000"
		res, err := testURL(url)
		if err != nil {
			t.Errorf("got %s, want nil", err)
		}
		if res != url {
			t.Errorf("got %s, want %s", res, url)
		}
	})

}
