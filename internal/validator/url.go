package validator

import (
	"net/url"
)

func Url(in string) error {
	_, err := url.ParseRequestURI(in)
	return err
}
