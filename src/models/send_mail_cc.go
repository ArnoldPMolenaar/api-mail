package models

type SendMailCc struct {
	SendMailID uint   `gorm:"primaryKey:true;not null;autoIncrement:false"`
	Cc         string `gorm:"primaryKey:true;not null;autoIncrement:false"`

	// Relationships.
	SendMail SendMail `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:SendMailID;references:ID"`
}
