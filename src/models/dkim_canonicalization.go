package models

type DkimCanonicalization struct {
	Name string `gorm:"primaryKey:true;not null;autoIncrement:false"`
}
