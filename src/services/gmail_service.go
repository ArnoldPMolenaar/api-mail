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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"os"
	"time"
)

// IsGmailAvailable checks if the gmail exists.
func IsGmailAvailable(app, mail string) (bool, error) {
	var count int64
	if result := database.Pg.Model(&models.Gmail{}).
		Joins("JOIN app_mails ON app_mails.id = gmails.app_mail_id").
		Where("app_mails.app_name = ? AND app_mails.mail_name = ?", app, mail).
		Count(&count); result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// IsGmailInCache checks if the gmail exists in the cache.
func IsGmailInCache(id uint) (bool, error) {
	key := gmailCacheKey(id)

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

// GetGmail gets the gmail.
func GetGmail(id uint, unscoped ...bool) (*models.Gmail, error) {
	gmail := &models.Gmail{}
	query := database.Pg

	if len(unscoped) > 0 && unscoped[0] {
		query = query.Unscoped()
	}

	if result := query.Preload("AppMail").Find(gmail, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	return gmail, nil
}

// GetGmailByID gets the gmail by ID.
func GetGmailByID(ID uint) (*models.Gmail, error) {
	gmail := &models.Gmail{}

	if result := database.Pg.Find(gmail, "id = ?", ID); result.Error != nil {
		return nil, result.Error
	}

	return gmail, nil
}

// GetGmailIDByAppMailID gets the gmail id by app mail id.
func GetGmailIDByAppMailID(appMailID uint) (uint, error) {
	var gmailID uint

	if result := database.Pg.Model(&models.Gmail{}).
		Where("app_mail_id = ?", appMailID).
		Pluck("id", &gmailID); result.Error != nil {
		return 0, result.Error
	}

	return gmailID, nil
}

// GetGmailFromCache gets the gmail from the cache.
func GetGmailFromCache(id uint) (*models.Gmail, error) {
	key := gmailCacheKey(id)

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(key).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var gmail models.Gmail
	if err := json.Unmarshal([]byte(value), &gmail); err != nil {
		return nil, err
	}

	return &gmail, nil
}

// CreateGmailOauthConfig creates a new oauth config.
func CreateGmailOauthConfig(clientID, secret string) *oauth2.Config {
	redirectUrl := fmt.Sprintf(
		"%sv1/oauth2/gmails/callback",
		os.Getenv("DOMAIN_NAME"),
	)

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Endpoint:     endpoints.Google,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send", "openid", "profile", "email"},
	}
}

// CreateGmail creates a new gmail.
func CreateGmail(req *requests.CreateGmail) (*models.Gmail, error) {
	gmailType := enums.Gmail
	gmail := &models.Gmail{
		ClientID: req.ClientID,
		Secret:   req.Secret,
		User:     req.User,
		AppMail: models.AppMail{
			AppName:  req.App,
			MailName: req.Mail,
		},
	}

	if appMail, err := GetAppMail(req.App, req.Mail); err != nil {
		return nil, err
	} else if appMail.ID != 0 {
		gmail.AppMail = appMail
	}

	if req.Primary {
		gmail.AppMail.PrimaryType = sql.NullString{String: *gmailType.ToString(), Valid: true}

		if gmail.AppMail.ID != 0 {
			if result := database.Pg.Save(gmail.AppMail); result.Error != nil {
				return nil, result.Error
			}
		}
	}

	if gmail.AppMail.ID == 0 {
		if result := database.Pg.FirstOrCreate(&models.Mail{
			Name: req.Mail,
		}); result.Error != nil {
			return nil, result.Error
		}
	}

	if result := database.Pg.Create(gmail); result.Error != nil {
		return nil, result.Error
	}

	return gmail, nil
}

// SetGmailToCache sets the gmail to the cache.
func SetGmailToCache(gmail *models.Gmail) error {
	key := gmailCacheKey(gmail.ID)

	value, err := json.Marshal(gmail)
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

// UpdateGmailToken updates an existing gmail token.
func UpdateGmailToken(gmail *models.Gmail, token *oauth2.Token) (*models.Gmail, error) {
	// Save the token into the Gmail record.
	gmail.AccessToken = sql.NullString{Valid: true, String: token.AccessToken}
	gmail.TokenType = sql.NullString{Valid: true, String: token.TokenType}
	gmail.Expiry = sql.NullTime{Valid: true, Time: token.Expiry}

	if token.RefreshToken != "" {
		gmail.RefreshToken = sql.NullString{Valid: true, String: token.RefreshToken}
	}

	if token.ExpiresIn != 0 {
		gmail.ExpiresIn = sql.NullInt64{Valid: true, Int64: token.ExpiresIn}
	} else {
		if expiresInFloat64, ok := token.Extra("expires_in").(float64); ok {
			gmail.ExpiresIn = sql.NullInt64{Valid: true, Int64: int64(expiresInFloat64)}
		}
	}

	// Update the Gmail record in the database.
	if result := database.Pg.Save(gmail); result.Error != nil {
		return nil, result.Error
	}

	return gmail, nil
}

// UpdateGmail updates a existing gmail.
func UpdateGmail(oldGmail *models.Gmail, req *requests.UpdateGmail) (*models.Gmail, error) {
	gmailType := enums.Gmail
	oldGmail.ClientID = req.ClientID
	oldGmail.Secret = req.Secret
	oldGmail.User = req.User

	if req.Primary && (!oldGmail.AppMail.PrimaryType.Valid || oldGmail.AppMail.PrimaryType.String != *gmailType.ToString()) {
		oldGmail.AppMail.PrimaryType = sql.NullString{String: *gmailType.ToString(), Valid: true}
	} else if !req.Primary && oldGmail.AppMail.PrimaryType.Valid && oldGmail.AppMail.PrimaryType.String == *gmailType.ToString() {
		oldGmail.AppMail.PrimaryType = sql.NullString{String: "", Valid: false}
	}

	if result := database.Pg.Save(oldGmail); result.Error != nil {
		return nil, result.Error
	}

	if result := database.Pg.Save(oldGmail.AppMail); result.Error != nil {
		return nil, result.Error
	}

	if isInCache, err := IsGmailInCache(oldGmail.ID); err != nil {
		return nil, err
	} else if isInCache {
		if err := SetGmailToCache(oldGmail); err != nil {
			return nil, err
		}
	}

	return oldGmail, nil
}

// DeleteGmail deletes a existing gmail.
func DeleteGmail(gmail *models.Gmail) error {
	if result := database.Pg.Delete(gmail); result.Error != nil {
		return result.Error
	}

	if isInCache, err := IsGmailInCache(gmail.ID); err != nil {
		return err
	} else if isInCache {
		if err := DeleteGmailFromCache(gmail.ID); err != nil {
			return err
		}
	}

	return nil
}

// DeleteGmailFromCache deletes a existing gmail from the cache.
func DeleteGmailFromCache(id uint) error {
	key := gmailCacheKey(id)

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(key).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// RestoreGmail restores a deleted gmail.
func RestoreGmail(gmail *models.Gmail) error {
	if result := database.Pg.Model(&gmail).Unscoped().Update("deleted_at", nil); result.Error != nil {
		return result.Error
	}

	return nil
}

// gmailCacheKey returns the key for the gmail cache.
func gmailCacheKey(id uint) string {
	return fmt.Sprintf("%s:%d", enums.Gmail, id)
}
