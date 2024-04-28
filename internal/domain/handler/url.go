package handler

import (
	"net/url"
)

type URL struct {
	Url string `json:"url"`
}

func (u URL) Validate() bool {
	if _, err := url.ParseRequestURI(u.Url); err != nil {
		return false
	}

	return true
}
