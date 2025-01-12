package models

// AppMailPrimaryType is an enum that contains Azure, Gmail or SMTP.
type AppMailPrimaryType struct {
	Name string `gorm:"primaryKey:true;not null;autoIncrement:false"`
}
