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
	"golang.org/x/oauth2/microsoft"
	"os"
	"time"
)

// IsAzureAvailable checks if the azure exists.
func IsAzureAvailable(app, mail string) (bool, error) {
	var count int64
	if result := database.Pg.Model(&models.Azure{}).
		Joins("JOIN app_mails ON app_mails.id = azures.app_mail_id").
		Where("app_mails.app_name = ? AND app_mails.mail_name = ?", app, mail).
		Count(&count); result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// IsAzureInCache checks if the azure exists in the cache.
func IsAzureInCache(id uint) (bool, error) {
	key := azureCacheKey(id)

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

// GetAzure gets the azure.
func GetAzure(id uint, unscoped ...bool) (*models.Azure, error) {
	azure := &models.Azure{}
	query := database.Pg

	if len(unscoped) > 0 && unscoped[0] {
		query = query.Unscoped()
	}

	if result := query.Preload("AppMail").Find(azure, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	return azure, nil
}

// GetAzureByID gets the azure by ID.
func GetAzureByID(ID uint) (*models.Azure, error) {
	azure := &models.Azure{}

	if result := database.Pg.Find(azure, "id = ?", ID); result.Error != nil {
		return nil, result.Error
	}

	return azure, nil
}

// GetAzureIDByAppMailID gets the azure id by app mail id.
func GetAzureIDByAppMailID(appMailID uint) (uint, error) {
	var azureID uint

	if result := database.Pg.Model(&models.Azure{}).
		Where("app_mail_id = ?", appMailID).
		Pluck("id", &azureID); result.Error != nil {
		return 0, result.Error
	}

	return azureID, nil
}

// GetAzureFromCache gets the azure from the cache.
func GetAzureFromCache(id uint) (*models.Azure, error) {
	key := azureCacheKey(id)

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(key).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var azure models.Azure
	if err := json.Unmarshal([]byte(value), &azure); err != nil {
		return nil, err
	}

	return &azure, nil
}

// CreateAzureOauthConfig creates a new oauth config.
func CreateAzureOauthConfig(clientID, tenantID, secret string) *oauth2.Config {
	redirectUrl := fmt.Sprintf(
		"%sv1/oauth2/azures/callback",
		os.Getenv("DOMAIN_NAME"),
	)

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Endpoint:     microsoft.AzureADEndpoint(tenantID),
		RedirectURL:  redirectUrl,
		Scopes:       []string{"openid", "offline_access", "user.read", "mail.send"},
	}
}

// CreateAzure creates a new azure.
func CreateAzure(req *requests.CreateAzure) (*models.Azure, error) {
	azureType := enums.Azure
	azure := &models.Azure{
		ClientID: req.ClientID,
		TenantID: req.TenantID,
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
		azure.AppMail = appMail
	}

	if req.Primary {
		azure.AppMail.PrimaryType = sql.NullString{String: *azureType.ToString(), Valid: true}

		if azure.AppMail.ID != 0 {
			if result := database.Pg.Save(azure.AppMail); result.Error != nil {
				return nil, result.Error
			}
		}
	}

	if azure.AppMail.ID == 0 {
		if result := database.Pg.FirstOrCreate(&models.Mail{
			Name: req.Mail,
		}); result.Error != nil {
			return nil, result.Error
		}
	}

	if result := database.Pg.Create(azure); result.Error != nil {
		return nil, result.Error
	}

	return azure, nil
}

// SetAzureToCache sets the azure to the cache.
func SetAzureToCache(azure *models.Azure) error {
	key := azureCacheKey(azure.ID)

	value, err := json.Marshal(azure)
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

// UpdateAzureToken updates an existing azure token.
func UpdateAzureToken(azure *models.Azure, token *oauth2.Token) (*models.Azure, error) {
	// Save the token into the Azure record.
	azure.AccessToken = sql.NullString{Valid: true, String: token.AccessToken}
	azure.TokenType = sql.NullString{Valid: true, String: token.TokenType}
	azure.Expiry = sql.NullTime{Valid: true, Time: token.Expiry}

	if token.RefreshToken != "" {
		azure.RefreshToken = sql.NullString{Valid: true, String: token.RefreshToken}
	}

	if token.ExpiresIn != 0 {
		azure.ExpiresIn = sql.NullInt64{Valid: true, Int64: token.ExpiresIn}
	} else {
		if expiresInFloat64, ok := token.Extra("expires_in").(float64); ok {
			azure.ExpiresIn = sql.NullInt64{Valid: true, Int64: int64(expiresInFloat64)}
		}
	}

	// Update the Azure record in the database.
	if result := database.Pg.Save(azure); result.Error != nil {
		return nil, result.Error
	}

	return azure, nil
}

// UpdateAzure updates a existing azure.
func UpdateAzure(oldAzure *models.Azure, req *requests.UpdateAzure) (*models.Azure, error) {
	azureType := enums.Azure
	oldAzure.ClientID = req.ClientID
	oldAzure.TenantID = req.TenantID
	oldAzure.Secret = req.Secret
	oldAzure.User = req.User

	if req.Primary && (!oldAzure.AppMail.PrimaryType.Valid || oldAzure.AppMail.PrimaryType.String != *azureType.ToString()) {
		oldAzure.AppMail.PrimaryType = sql.NullString{String: *azureType.ToString(), Valid: true}
	} else if !req.Primary && oldAzure.AppMail.PrimaryType.Valid && oldAzure.AppMail.PrimaryType.String == *azureType.ToString() {
		oldAzure.AppMail.PrimaryType = sql.NullString{String: "", Valid: false}
	}

	if result := database.Pg.Save(oldAzure); result.Error != nil {
		return nil, result.Error
	}

	if result := database.Pg.Save(oldAzure.AppMail); result.Error != nil {
		return nil, result.Error
	}

	if isInCache, err := IsAzureInCache(oldAzure.ID); err != nil {
		return nil, err
	} else if isInCache {
		if err := SetAzureToCache(oldAzure); err != nil {
			return nil, err
		}
	}

	return oldAzure, nil
}

// DeleteAzure deletes a existing azure.
func DeleteAzure(azure *models.Azure) error {
	if result := database.Pg.Delete(azure); result.Error != nil {
		return result.Error
	}

	if isInCache, err := IsAzureInCache(azure.ID); err != nil {
		return err
	} else if isInCache {
		if err := DeleteAzureFromCache(azure.ID); err != nil {
			return err
		}
	}

	return nil
}

// DeleteAzureFromCache deletes a existing azure from the cache.
func DeleteAzureFromCache(id uint) error {
	key := azureCacheKey(id)

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(key).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// RestoreAzure restores a deleted azure.
func RestoreAzure(azure *models.Azure) error {
	if result := database.Pg.Model(&azure).Unscoped().Update("deleted_at", nil); result.Error != nil {
		return result.Error
	}

	return nil
}

// azureCacheKey returns the key for the azure cache.
func azureCacheKey(id uint) string {
	return fmt.Sprintf("%s:%d", enums.Azure, id)
}
