package services

import (
	"api-mail/main/src/cache"
	"api-mail/main/src/database"
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/enums"
	"api-mail/main/src/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/valkey-io/valkey-go"
	"os"
	"time"
)

// IsSmtpAvailable checks if the smtp exists.
func IsSmtpAvailable(app, mail string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.Smtp{}, "app_name = ? AND mail_name = ?", app, mail); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// IsSmtpInCache checks if the smtp exists in the cache.
func IsSmtpInCache(app, mail string) (bool, error) {
	key := cacheKey(app, mail)

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(key).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// GetSmtp gets the smtp.
func GetSmtp(app, mail string, unscoped ...bool) (*models.Smtp, error) {
	smtp := &models.Smtp{}
	query := database.Pg

	if len(unscoped) > 0 && unscoped[0] {
		query = query.Unscoped()
	}

	if result := query.Find(smtp, "app_name = ? AND mail_name = ?", app, mail); result.Error != nil {
		return nil, result.Error
	}

	return smtp, nil
}

// GetSmtpFromCache gets the smtp from the cache.
func GetSmtpFromCache(app, mail string) (*models.Smtp, error) {
	key := cacheKey(app, mail)

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(key).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var smtp models.Smtp
	if err := json.Unmarshal([]byte(value), &smtp); err != nil {
		return nil, err
	}

	return &smtp, nil
}

// CreateSmtp creates a new smtp.
func CreateSmtp(req *requests.CreateSmtp) error {
	smtpType := enums.SMTP
	smtp := &models.Smtp{
		AppName:                  req.App,
		MailName:                 req.Mail,
		Username:                 req.Username,
		Password:                 req.Password,
		Host:                     req.Host,
		Port:                     req.Port,
		DkimPrivateKey:           req.DkimPrivateKey,
		DkimDomain:               req.DkimDomain,
		DkimCanonicalizationName: req.DkimCanonicalization,
	}

	if err := smtp.EncryptPassword(); err != nil {
		return err
	}

	if result := database.Pg.FirstOrCreate(&models.Mail{
		Name: req.Mail,
	}); result.Error != nil {
		return result.Error
	}

	if result := database.Pg.Create(smtp); result.Error != nil {
		return result.Error
	}

	if result := database.Pg.Create(&models.AppMail{
		AppName:  req.App,
		MailName: req.Mail,
		MailType: smtpType.ToString(),
		Primary:  req.Primary,
	}); result.Error != nil {
		return result.Error
	}

	return nil
}

// SetSmtpToCache sets the smtp to the cache.
func SetSmtpToCache(smtp *models.Smtp) error {
	key := cacheKey(smtp.AppName, smtp.MailName)

	value, err := json.Marshal(smtp)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(key).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// UpdateSmtp updates a existing smtp.
func UpdateSmtp(oldSmtp *models.Smtp, req *requests.UpdateSmtp) error {
	smtpType := enums.SMTP
	smtp := &models.Smtp{
		AppName:                  oldSmtp.AppName,
		MailName:                 oldSmtp.MailName,
		Username:                 req.Username,
		Password:                 oldSmtp.Password,
		Host:                     req.Host,
		Port:                     req.Port,
		DkimPrivateKey:           req.DkimPrivateKey,
		DkimDomain:               req.DkimDomain,
		DkimCanonicalizationName: req.DkimCanonicalization,
		CreatedAt:                oldSmtp.CreatedAt,
	}

	if req.Password != "" {
		smtp.Password = req.Password
		if err := smtp.EncryptPassword(); err != nil {
			return err
		}
	}

	if result := database.Pg.Save(smtp); result.Error != nil {
		return result.Error
	}

	if result := database.Pg.Save(&models.AppMail{
		AppName:  oldSmtp.AppName,
		MailName: oldSmtp.MailName,
		MailType: smtpType.ToString(),
		Primary:  req.Primary,
	}); result.Error != nil {
		return result.Error
	}

	if isInCache, err := IsSmtpInCache(oldSmtp.AppName, oldSmtp.MailName); err != nil {
		return err
	} else if isInCache {
		if err := SetSmtpToCache(smtp); err != nil {
			return err
		}
	}

	return nil
}

// DeleteSmtp deletes a existing smtp.
func DeleteSmtp(smtp *models.Smtp) error {
	if result := database.Pg.Delete(smtp); result.Error != nil {
		return result.Error
	}

	if isInCache, err := IsSmtpInCache(smtp.AppName, smtp.MailName); err != nil {
		return err
	} else if isInCache {
		if err := DeleteSmtpFromCache(smtp.AppName, smtp.MailName); err != nil {
			return err
		}
	}

	return nil
}

// DeleteSmtpFromCache deletes a existing smtp from the cache.
func DeleteSmtpFromCache(app, mail string) error {
	key := cacheKey(app, mail)

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(key).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// RestoreSmtp restores a deleted smtp.
func RestoreSmtp(smtp *models.Smtp) error {
	if result := database.Pg.Model(&smtp).Unscoped().Update("deleted_at", nil); result.Error != nil {
		return result.Error
	}

	return nil
}

// cacheKey returns the key for the smtp cache.
func cacheKey(app, mail string) string {
	return fmt.Sprintf("%s:%s:%s", enums.SMTP, app, mail)
}
