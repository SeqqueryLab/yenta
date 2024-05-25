package model

import "testing"

func TestBindingSimple(t *testing.T) {
	b := Routing{"simple", "*.orange.*", "more.complicated.key"}

	t.Run("simple or composite binding", func(t *testing.T) {
		got, err := b.Simple()
		want := true
		if want != got || err != nil {
			t.Errorf("want %t, got %t. error: %s", want, got, err)
		}
	})
}
