package models

import (
	"api-mail/main/src/utils"
	"gorm.io/gorm"
	"os"
)

type Smtp struct {
	gorm.Model
	AppMailID                uint   `gorm:"not null"`
	Username                 string `gorm:"not null"`
	Password                 string `gorm:"not null"`
	Host                     string `gorm:"not null"`
	Port                     int    `gorm:"not null"`
	DkimPrivateKey           *string
	DkimDomain               *string
	DkimCanonicalizationName *string

	// Relationships.
	AppMail              AppMail               `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppMailID;references:ID"`
	DkimCanonicalization *DkimCanonicalization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:DkimCanonicalizationName;references:Name"`
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
