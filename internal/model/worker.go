package model

type Worker func(map[string]interface{}) map[string]interface{}

func (w Worker) Do(ch chan interface{}, args map[string]interface{}) {
	go func() {
		ch <- w(args)
	}()
}
