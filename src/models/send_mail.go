package models

import (
	"time"
)

type SendMail struct {
	ID        uint   `gorm:"primarykey"`
	AppName   string `gorm:"not null"`
	MailName  string `gorm:"not null"`
	MailType  string `gorm:"not null"`
	FromName  string
	FromMail  string
	To        string    `gorm:"not null"`
	Subject   string    `gorm:"not null"`
	Body      string    `gorm:"not null"`
	MimeType  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`

	// Relationships.
	App  App           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppName;references:Name"`
	Mail Mail          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MailName;references:Name"`
	Type MailType      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MailType;references:Name"`
	Ccs  []SendMailCc  `gorm:"foreignKey:SendMailID"`
	Bccs []SendMailBcc `gorm:"foreignKey:SendMailID"`
}
