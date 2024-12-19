package models

type App struct {
	Name string `gorm:"primaryKey:true;not null;autoIncrement:false"`
}
