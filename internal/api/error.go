package api

import "fmt"

type Error chan error

func (e Error) Error() string {
	defer close(e)
	select {
	case err := <-e:
		return fmt.Sprintf("%s", err)
	}
}
