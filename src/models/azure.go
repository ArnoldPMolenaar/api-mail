package models

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type Azure struct {
	gorm.Model
	AppMailID uint   `gorm:"not null"`
	ClientID  string `gorm:"primaryKey:true;not null;autoIncrement:false"`
	TenantID  string `gorm:"not null"`
	Secret    string `gorm:"not null"`
	Token     string
	User      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`

	// Relationships.
	AppMail AppMail `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppMailID;references:ID"`
}
