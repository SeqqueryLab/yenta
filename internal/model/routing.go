package model

import (
	"errors"
	"regexp"
)

type Routing []string

func (b *Routing) Simple() (bool, error) {
	if len(*b) == 0 {
		return true, errors.New("binding keys are not found")
	}
	re := regexp.MustCompile(`\.+`)
	var res bool
	for _, val := range *b {
		res = res || re.MatchString(val)
	}
	return res, nil
}
