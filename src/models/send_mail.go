package models

import (
	"time"
)

type SendMail struct {
	ID        uint   `gorm:"primarykey"`
	From      string `gorm:"not null"`
	To        string `gorm:"not null"`
	Subject   string `gorm:"not null"`
	Body      string `gorm:"not null"`
	MimeType  string `gorm:"not null"`
	CreatedAt time.Time
}
