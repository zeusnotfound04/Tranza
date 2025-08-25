package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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

// Generate creates a universal API key that works with everything
func (s *APIKeyService) Generate(ctx context.Context, userID uuid.UUID, label string, ttlHours int) (string, error) {
	// Universal scopes that allow access to all features
	universalScopes := []string{
		"*", // Wildcard scope for all permissions
	}

	return s.GenerateWithScopes(ctx, userID, label, ttlHours, universalScopes, "universal", "", "")
}

// GenerateBotKey creates a universal API key (same as Generate for now)
// Keeping this for backward compatibility but it now creates the same universal key
func (s *APIKeyService) GenerateBotKey(ctx context.Context, userID uuid.UUID, label string, workspaceID, botUserID string, ttlHours int) (string, error) {
	// For now, bot keys are the same as universal keys
	// In the future, you can differentiate them if needed
	universalScopes := []string{
		"*", // Wildcard scope for all permissions
	}

	return s.GenerateWithScopes(ctx, userID, label, ttlHours, universalScopes, "universal", workspaceID, botUserID)
}

// GenerateWithScopes creates an API key with specific scopes and type
func (s *APIKeyService) GenerateWithScopes(ctx context.Context, userID uuid.UUID, label string, ttlHours int, scopes []string, keyType, workspaceID, botUserID string) (string, error) {
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
		UserID:       userID,
		KeyHash:      hash,
		Label:        label,
		KeyType:      keyType,
		IsActive:     true,
		ExpiresAt:    expiresAt,
		BotWorkspace: workspaceID,
		BotUserID:    botUserID,
		RateLimit:    1000, // Default rate limit
	}

	// Set scopes
	if err := key.SetScopes(scopes); err != nil {
		return "", fmt.Errorf("failed to set scopes: %w", err)
	}

	if err := s.Repo.Create(ctx, key); err != nil {
		return "", fmt.Errorf("failed to create API key: %w", err)
	}

	return rawKey, nil
}

// Validate validates an API key and returns it if valid
func (s *APIKeyService) Validate(ctx context.Context, rawKey string) (*models.APIKey, error) {
	hash := utils.HashKey(rawKey)
	key, err := s.Repo.FindByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	if key == nil || !key.CanMakeRequest() {
		return nil, errors.New("invalid or expired API key")
	}

	// Increment usage counter
	key.IncrementUsage()
	_ = s.Repo.UpdateUsage(ctx, key.ID)

	return key, nil
}

// ValidateWithScope validates an API key and checks if it has the required scope
func (s *APIKeyService) ValidateWithScope(ctx context.Context, rawKey string, requiredScope string) (*models.APIKey, error) {
	key, err := s.Validate(ctx, rawKey)
	if err != nil {
		return nil, err
	}

	if !key.HasScope(requiredScope) {
		return nil, errors.New("insufficient permissions for this operation")
	}

	return key, nil
}

// RotateKey generates a new key for an existing API key entry
func (s *APIKeyService) RotateKey(ctx context.Context, keyID uint, userID uuid.UUID) (string, error) {
	// Get existing key
	existingKey, err := s.Repo.FindByID(ctx, keyID)
	if err != nil {
		return "", err
	}

	if existingKey.UserID != userID {
		return "", errors.New("unauthorized")
	}

	// Generate new raw key
	rawKey, err := utils.GenerateSecureKey()
	if err != nil {
		return "", err
	}

	// Update the key hash
	newHash := utils.HashKey(rawKey)
	err = s.Repo.UpdateKeyHash(ctx, keyID, newHash)
	if err != nil {
		return "", fmt.Errorf("failed to rotate key: %w", err)
	}

	return rawKey, nil
}

// GetUsageStats returns usage statistics for an API key
func (s *APIKeyService) GetUsageStats(ctx context.Context, keyID uint, userID uuid.UUID) (*APIKeyUsageStats, error) {
	key, err := s.Repo.FindByID(ctx, keyID)
	if err != nil {
		return nil, err
	}

	if key.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return &APIKeyUsageStats{
		KeyID:      key.ID,
		Label:      key.Label,
		KeyType:    key.KeyType,
		UsageCount: key.UsageCount,
		RateLimit:  key.RateLimit,
		LastUsedAt: key.LastUsedAt,
		CreatedAt:  key.CreatedAt,
		Scopes:     key.GetScopes(),
	}, nil
}

// Revoke revokes an API key
func (s *APIKeyService) Revoke(ctx context.Context, keyID uint, userID uuid.UUID) error {
	return s.Repo.RevokeByID(ctx, keyID, userID)
}

// APIKeyUsageStats represents usage statistics for an API key
type APIKeyUsageStats struct {
	KeyID      uint      `json:"key_id"`
	Label      string    `json:"label"`
	KeyType    string    `json:"key_type"`
	UsageCount int64     `json:"usage_count"`
	RateLimit  int       `json:"rate_limit"`
	LastUsedAt time.Time `json:"last_used_at"`
	CreatedAt  time.Time `json:"created_at"`
	Scopes     []string  `json:"scopes"`
}
