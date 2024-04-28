package handler

import (
	"errors"
	"net/url"
)

type URL struct {
	Url string `json:"url"`
}

func (u URL) Validate() error {
	if _, err := url.ParseRequestURI(u.Url); err != nil {
		return errors.New("invalid url")
	}

	return nil
}
