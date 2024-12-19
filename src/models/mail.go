package models

type Mail struct {
	Name string `gorm:"primaryKey:true;not null;autoIncrement:false"`
}
