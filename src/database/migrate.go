package database

import (
	"api-mail/main/src/models"
	"gorm.io/gorm"
)

// Migrate the database schema.
// See: https://gorm.io/docs/migration.html#Auto-Migration
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		models.App{},
		models.Mail{},
		models.MailType{},
		models.Smtp{},
		models.Azure{},
		models.Gmail{},
		models.AppMail{},
		models.DkimCanonicalization{},
		models.SendMail{},
		models.SendMailCc{},
		models.SendMailBcc{})
	if err != nil {
		return err
	}

	// Seed App.
	apps := []string{"Admin"}
	for _, app := range apps {
		if err := db.FirstOrCreate(&models.App{}, models.App{Name: app}).Error; err != nil {
			return err
		}
	}

	// Seed MailType.
	mailTypes := []string{"Azure", "Gmail", "SMTP"}
	for _, mailType := range mailTypes {
		if err := db.FirstOrCreate(&models.MailType{}, models.MailType{Name: mailType}).Error; err != nil {
			return err
		}
	}

	// Seed DkimCanonicalization.
	dkimCanonicalization := []string{"Simple", "Relaxed"}
	for _, dkimCanonicalization := range dkimCanonicalization {
		if err := db.FirstOrCreate(&models.DkimCanonicalization{}, models.DkimCanonicalization{Name: dkimCanonicalization}).Error; err != nil {
			return err
		}
	}

	return nil
}
