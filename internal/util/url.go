package util

import (
	"errors"
	"regexp"
)

func testURL(url string) (string, error) {
	if ok := regexp.MustCompile(`^amqp://\w+:.+@[a-z.]+:\d+$`).MatchString(url); !ok {
		return "", errors.New("invalid RabbitMQ url")
	}
	return url, nil
}
