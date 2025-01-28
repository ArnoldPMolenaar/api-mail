package models

import (
	"time"
)

type SendMailAttachment struct {
	ID         uint      `gorm:"primarykey"`
	SendMailID uint      `gorm:"not null"`
	FileName   string    `gorm:"not null"`
	FileType   string    `gorm:"not null"`
	FileSize   int64     `gorm:"not null"`
	FileData   []byte    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`

	// Relationships.
	SendMail SendMail `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:SendMailID;references:ID"`
}
