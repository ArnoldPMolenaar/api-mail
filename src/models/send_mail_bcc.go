package models

type SendMailBcc struct {
	SendMailID uint   `gorm:"primaryKey:true;not null;autoIncrement:false"`
	Bcc        string `gorm:"primaryKey:true;not null;autoIncrement:false"`

	// Relationships.
	SendMail SendMail `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:SendMailID;references:ID"`
}
