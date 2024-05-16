package util

import (
	"reflect"
	"testing"
)

func TestFibonacci(t *testing.T) {
	t.Run("return first element of fibonacci sequence", func(t *testing.T) {
		got, _ := Fibonacci(1)
		want := []int{0}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v\n", got, want)
		}
	})

	t.Run("return first two elements of fibonacci sequence", func(t *testing.T) {
		got, _ := Fibonacci(2)
		want := []int{0, 1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v\n", got, want)
		}
	})

	t.Run("return first five elements of fibonacci sequence", func(t *testing.T) {
		got, _ := Fibonacci(5)
		want := []int{0, 1, 1, 2, 3}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v\n", got, want)
		}
	})

	t.Run("return first ten elements of fibonacci sequence", func(t *testing.T) {
		got, _ := Fibonacci(10)
		want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v\n", got, want)
		}
	})
}
