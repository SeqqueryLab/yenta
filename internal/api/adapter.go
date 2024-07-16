package api

// service
type Service struct {
	url        string
	publishers map[string]chan interface{}
	err        *Error
}
