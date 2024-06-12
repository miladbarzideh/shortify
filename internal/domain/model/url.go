package model

import (
	"net/url"
	"time"
)

type URL struct {
	ID        uint `gorm:"primaryKey; auto_increment"`
	LongURL   string
	ShortCode string `gorm:"unique; size:20; index'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type URLData struct {
	URL string `json:"url"`
}

func (u URLData) Validate() bool {
	parsedURL, err := url.Parse(u.URL)
	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		return true
	}

	return false
}
