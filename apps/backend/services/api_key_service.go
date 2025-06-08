package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
	"github.com/zeusnotfound04/Tranza/utils"
)

type APIKeyService struct {
	Repo *repositories.APIKeyRepository
}

func NewAPIKeyService(repo *repositories.APIKeyRepository) *APIKeyService {
	return &APIKeyService{Repo: repo}
}

func (s *APIKeyService) Generate(ctx context.Context, userID uint, label string, ttlHours int) (string, error) {
	rawKey, err := utils.GenerateSecureKey()
	if err != nil {
		return "", err
	}

	hash := utils.HashKey(rawKey)
	var expiresAt *time.Time
	if ttlHours > 0 {
		exp := time.Now().Add(time.Duration(ttlHours) * time.Hour)
		expiresAt = &exp
	}

	key := &models.APIKey{
		UserID:    userID,
		KeyHash:   hash,
		Label:     label,
		IsActive:  true,
		ExpiresAt: expiresAt,
	}

	if err := s.Repo.Create(ctx, key); err != nil {
		return "", fmt.Errorf("failed to create API key: %w", err)
	}

	return rawKey, nil
}

func (s *APIKeyService) Validate(ctx context.Context, rawKey string) (*models.APIKey, error) {
	hash := utils.HashKey(rawKey)
	key, err := s.Repo.FindByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	if key == nil || !key.IsActive || key.IsExpired() {
		return nil, errors.New("invalid or expired API key")
	}

	_ = s.Repo.UpdateUsage(ctx, key.ID)
	return key, nil
}

func (s *APIKeyService) Revoke(ctx context.Context, keyID, userID uint) error {
	return s.Repo.RevokeByID(ctx, keyID, userID)
}