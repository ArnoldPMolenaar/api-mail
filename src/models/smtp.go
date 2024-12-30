package models

import (
	"api-mail/main/src/utils"
	"database/sql"
	"os"
	"time"
)

type Smtp struct {
	AppName                  string `gorm:"primaryKey:true;not null;autoIncrement:false"`
	MailName                 string `gorm:"primaryKey:true;not null;autoIncrement:false"`
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
	Mail                 Mail                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MailName;references:Name"`
	DkimCanonicalization DkimCanonicalization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:DkimCanonicalizationName;references:Name"`
}

// EncryptPassword encrypts the SMTP password.
func (s *Smtp) EncryptPassword() error {
	key := os.Getenv("PASSWORD_ENCRYPTION_KEY")
	encryptedPassword, err := utils.Encrypt(key, s.Password)

	if err != nil {
		return err
	}

	s.Password = encryptedPassword

	return nil
}

// DecryptPassword decrypts the SMTP password.
func (s *Smtp) DecryptPassword() (string, error) {
	key := os.Getenv("PASSWORD_ENCRYPTION_KEY")

	return utils.Decrypt(key, s.Password)
}
