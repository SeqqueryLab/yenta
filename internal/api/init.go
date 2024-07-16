package api

func New(url string) (*Service, error) {
	_, err := util.testURL(url)
	if err != nil {
		return &Service{}, err
	}
	return &Service{url: url, publishers: make(map[string]chan interface{})}, nil
}
