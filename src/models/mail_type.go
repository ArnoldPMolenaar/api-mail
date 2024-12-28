package models

// MailType is an enum that contains Azure, Gmail or SMTP.
type MailType struct {
	Name string `gorm:"primaryKey:true;not null;autoIncrement:false"`
}
