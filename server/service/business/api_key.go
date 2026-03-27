package business

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
)

type APIKeyService struct{}

func (s *APIKeyService) GenerateAPIKey(name string, createdBy uint, expiresAt *time.Time) (*business.APIKey, string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, "", err
	}
	rawKey := hex.EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])

	apiKey := &business.APIKey{
		Name:      name,
		KeyHash:   keyHash,
		IsActive:  true,
		ExpiresAt: expiresAt,
		CreatedBy: createdBy,
	}

	if err := global.GVA_DB.Create(apiKey).Error; err != nil {
		return nil, "", err
	}

	return apiKey, rawKey, nil
}

func (s *APIKeyService) GetAPIKey(id uint) (apiKey *business.APIKey, err error) {
	err = global.GVA_DB.First(&apiKey, id).Error
	return
}

func (s *APIKeyService) GetAPIKeys() (apiKeys []business.APIKey, err error) {
	err = global.GVA_DB.Order("created_at DESC").Find(&apiKeys).Error
	return
}

func (s *APIKeyService) ValidateAPIKey(rawKey string) (*business.APIKey, error) {
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])

	var apiKey business.APIKey
	if err := global.GVA_DB.Where("key_hash = ? AND is_active = ?", keyHash, true).First(&apiKey).Error; err != nil {
		return nil, err
	}

	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, nil
	}

	return &apiKey, nil
}

func (s *APIKeyService) RevokeAPIKey(id uint) error {
	return global.GVA_DB.Model(&business.APIKey{}).Where("id = ?", id).Update("is_active", false).Error
}

func (s *APIKeyService) UpdateLastUsed(id uint) error {
	now := time.Now()
	return global.GVA_DB.Model(&business.APIKey{}).Where("id = ?", id).Update("last_used_at", &now).Error
}

func (s *APIKeyService) DeleteAPIKey(id uint) error {
	return global.GVA_DB.Delete(&business.APIKey{}, id).Error
}
