package yenta

import "fmt"

type YentaError chan error

func (err YentaError) Error() string {

	return fmt.Sprintf("yenta error: %s", <-err)
}
