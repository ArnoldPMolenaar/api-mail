package models

import (
	"database/sql"
	"time"
)

type Azure struct {
	AppName   string `gorm:"primaryKey:true;not null;autoIncrement:false"`
	ClientID  string `gorm:"primaryKey:true;not null;autoIncrement:false"`
	TenantID  string `gorm:"not null"`
	Secret    string `gorm:"not null"`
	Token     string
	User      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`

	// Relationships.
	App App `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppName;references:Name"`
}