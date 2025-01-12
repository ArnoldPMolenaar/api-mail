package models

import "database/sql"

type AppMail struct {
	ID          uint   `gorm:"primaryKey"`
	AppName     string `gorm:"not null;index:idx_app_mail,unique,priority:1"`
	MailName    string `gorm:"not null;index:idx_app_mail,unique,priority:2"`
	PrimaryType sql.NullString

	// Relationships.
	App   App                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppName;references:Name"`
	Mail  Mail                `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MailName;references:Name"`
	Type  *AppMailPrimaryType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PrimaryType;references:Name"`
	Smtp  *Smtp
	Gmail *Gmail
	Azure *Azure
}
