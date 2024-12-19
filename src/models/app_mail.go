package models

type AppMail struct {
	AppName  string `gorm:"primaryKey:true;not null;autoIncrement:false"`
	MailName string `gorm:"primaryKey:true;not null;autoIncrement:false"`
	Primary  bool   `gorm:"not null;default:false"`

	// Relationships.
	App  App  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppName;references:Name"`
	Mail Mail `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:RoleName;references:Name"`
}
