package services

import (
	"api-mail/main/src/cache"
	"api-mail/main/src/database"
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/enums"
	"api-mail/main/src/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/valkey-io/valkey-go"
	"os"
	"time"
)

// IsSmtpAvailable checks if the smtp exists.
func IsSmtpAvailable(app, mail string) (bool, error) {
	var count int64
	if result := database.Pg.Model(&models.Smtp{}).
		Joins("JOIN app_mails ON app_mails.id = smtps.app_mail_id").
		Where("app_mails.app_name = ? AND app_mails.mail_name = ?", app, mail).
		Count(&count); result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// IsSmtpInCache checks if the smtp exists in the cache.
func IsSmtpInCache(id uint) (bool, error) {
	key := smtpCacheKey(id)

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

// GetSmtpIDByAppMailID gets the smtp id by app mail id.
func GetSmtpIDByAppMailID(appMailID uint) (uint, error) {
	var smtpID uint

	if result := database.Pg.Model(&models.Smtp{}).
		Where("app_mail_id = ?", appMailID).
		Pluck("id", &smtpID); result.Error != nil {
		return 0, result.Error
	}

	return smtpID, nil
}

// GetSmtp gets the smtp.
func GetSmtp(id uint, unscoped ...bool) (*models.Smtp, error) {
	smtp := &models.Smtp{}
	query := database.Pg

	if len(unscoped) > 0 && unscoped[0] {
		query = query.Unscoped()
	}

	if result := query.Preload("AppMail").Find(smtp, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	return smtp, nil
}

// GetSmtpFromCache gets the smtp from the cache.
func GetSmtpFromCache(id uint) (*models.Smtp, error) {
	key := smtpCacheKey(id)

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
func CreateSmtp(req *requests.CreateSmtp) (*models.Smtp, error) {
	smtpType := enums.SMTP
	smtp := &models.Smtp{
		Username:                 req.Username,
		Password:                 req.Password,
		Host:                     req.Host,
		Port:                     req.Port,
		DkimPrivateKey:           req.DkimPrivateKey,
		DkimDomain:               req.DkimDomain,
		DkimCanonicalizationName: req.DkimCanonicalization,
		AppMail: models.AppMail{
			AppName:  req.App,
			MailName: req.Mail,
		},
	}

	if appMail, err := GetAppMail(req.App, req.Mail); err != nil {
		return nil, err
	} else if appMail.ID != 0 {
		smtp.AppMail = appMail
	}

	if req.Primary {
		smtp.AppMail.PrimaryType = sql.NullString{String: *smtpType.ToString(), Valid: true}

		if smtp.AppMail.ID != 0 {
			if result := database.Pg.Save(smtp.AppMail); result.Error != nil {
				return nil, result.Error
			}
		}
	}

	if err := smtp.EncryptPassword(); err != nil {
		return nil, err
	}

	if smtp.AppMail.ID == 0 {
		if result := database.Pg.FirstOrCreate(&models.Mail{
			Name: req.Mail,
		}); result.Error != nil {
			return nil, result.Error
		}
	}

	if result := database.Pg.Create(smtp); result.Error != nil {
		return nil, result.Error
	}

	return smtp, nil
}

// SetSmtpToCache sets the smtp to the cache.
func SetSmtpToCache(smtp *models.Smtp) error {
	key := smtpCacheKey(smtp.ID)

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
func UpdateSmtp(oldSmtp *models.Smtp, req *requests.UpdateSmtp) (*models.Smtp, error) {
	smtpType := enums.SMTP
	oldSmtp.Username = req.Username
	oldSmtp.Host = req.Host
	oldSmtp.Port = req.Port
	oldSmtp.DkimPrivateKey = req.DkimPrivateKey
	oldSmtp.DkimDomain = req.DkimDomain
	oldSmtp.DkimCanonicalizationName = req.DkimCanonicalization

	if req.Password != "" {
		oldSmtp.Password = req.Password
		if err := oldSmtp.EncryptPassword(); err != nil {
			return nil, err
		}
	}

	if req.Primary && (!oldSmtp.AppMail.PrimaryType.Valid || oldSmtp.AppMail.PrimaryType.String != *smtpType.ToString()) {
		oldSmtp.AppMail.PrimaryType = sql.NullString{String: *smtpType.ToString(), Valid: true}
	} else if !req.Primary && oldSmtp.AppMail.PrimaryType.Valid && oldSmtp.AppMail.PrimaryType.String == *smtpType.ToString() {
		oldSmtp.AppMail.PrimaryType = sql.NullString{String: "", Valid: false}
	}

	if result := database.Pg.Save(oldSmtp); result.Error != nil {
		return nil, result.Error
	}

	if result := database.Pg.Save(oldSmtp.AppMail); result.Error != nil {
		return nil, result.Error
	}

	if isInCache, err := IsSmtpInCache(oldSmtp.ID); err != nil {
		return nil, err
	} else if isInCache {
		if err := SetSmtpToCache(oldSmtp); err != nil {
			return nil, err
		}
	}

	return oldSmtp, nil
}

// DeleteSmtp deletes a existing smtp.
func DeleteSmtp(smtp *models.Smtp) error {
	if result := database.Pg.Delete(smtp); result.Error != nil {
		return result.Error
	}

	if isInCache, err := IsSmtpInCache(smtp.ID); err != nil {
		return err
	} else if isInCache {
		if err := DeleteSmtpFromCache(smtp.ID); err != nil {
			return err
		}
	}

	return nil
}

// DeleteSmtpFromCache deletes an existing smtp from the cache.
func DeleteSmtpFromCache(id uint) error {
	key := smtpCacheKey(id)

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

// smtpCacheKey returns the key for the smtp cache.
func smtpCacheKey(id uint) string {
	return fmt.Sprintf("%s:%d", enums.SMTP, id)
}
