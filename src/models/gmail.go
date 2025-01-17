package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type Gmail struct {
	gorm.Model
	AppMailID    uint   `gorm:"not null"`
	ClientID     string `gorm:"not null"`
	Secret       string `gorm:"not null"`
	AccessToken  sql.NullString
	RefreshToken sql.NullString
	TokenType    sql.NullString
	Expiry       sql.NullTime
	ExpiresIn    sql.NullInt64
	User         string `gorm:"not null"`

	// Relationships.
	AppMail AppMail `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppMailID;references:ID"`
}
