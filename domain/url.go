package domain

import (
	"time"
)

type URL struct {
	ID        uint `gorm:"primaryKey; auto_increment"`
	LongURL   string
	ShortCode string `gorm:"unique; size:20; index'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
