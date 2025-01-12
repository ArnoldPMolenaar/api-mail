package models

import (
	"time"
)

type SendMail struct {
	ID          uint   `gorm:"primarykey"`
	AppMailID   uint   `gorm:"not null"`
	PrimaryType string `gorm:"not null"`
	FromName    string
	FromMail    string
	To          string    `gorm:"not null"`
	Subject     string    `gorm:"not null"`
	Body        string    `gorm:"not null"`
	MimeType    string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`

	// Relationships.
	AppMail AppMail            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppMailID;references:ID"`
	Type    AppMailPrimaryType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PrimaryType;references:Name"`
	Ccs     []SendMailCc       `gorm:"foreignKey:SendMailID"`
	Bccs    []SendMailBcc      `gorm:"foreignKey:SendMailID"`
}
