package models

import (
	"gorm.io/gorm"
	"time"
)

type Gmail struct {
	gorm.Model
	AppMailID    uint   `gorm:"not null"`
	ClientID     string `gorm:"not null"`
	Secret       string `gorm:"not null"`
	AccessToken  *string
	RefreshToken *string
	TokenType    *string
	Expiry       *time.Time
	ExpiresIn    *int64
	User         string `gorm:"not null"`

	// Relationships.
	AppMail AppMail `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppMailID;references:ID"`
}
