package models

import (
	"database/sql"
	"time"
)

type Smtp struct {
	AppName                  string `gorm:"primaryKey:true;not null;autoIncrement:false"`
	User                     string `gorm:"primaryKey:true;not null;autoIncrement:false"`
	Password                 string `gorm:"not null"`
	Host                     string `gorm:"not null"`
	Port                     int    `gorm:"not null"`
	DkimPrivateKey           string
	DkimDomain               string
	DkimCanonicalizationName string
	CreatedAt                time.Time
	UpdatedAt                time.Time
	DeletedAt                sql.NullTime `gorm:"index"`

	// Relationships.
	App                  App                  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppName;references:Name"`
	DkimCanonicalization DkimCanonicalization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:DkimCanonicalizationName;references:Name"`
}
