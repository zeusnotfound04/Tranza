package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
	"github.com/zeusnotfound04/Tranza/utils"
)

type AuthService struct {
	userRepo      repositories.UserRepository
	jwtService    utils.JWTService
	oauthService  OAuthService
	walletService *WalletService
}

func NewAuthService(userRepo repositories.UserRepository, jwtService utils.JWTService, oauthService OAuthService, walletService *WalletService) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		jwtService:    jwtService,
		oauthService:  oauthService,
		walletService: walletService,
	}
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !s.verifyPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user exists
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		Provider:  "local",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *AuthService) GetAuthURL(provider, state string) (string, error) {
	return s.oauthService.GetAuthURL(provider, state, "")
}

func (s *AuthService) HandleOAuthCallback(ctx context.Context, req models.OAuthCallbackRequest) (*models.AuthResponse, error) {
	// Exchange code for user info
	oauthUser, err := s.oauthService.ExchangeCodeForUser(ctx, req.Provider, req.Code, req.RedirectURI)
	if err != nil {
		return nil, err
	}

	// Find or create user
	user, err := s.userRepo.FindByProviderID(ctx, req.Provider, oauthUser.ID)
	if err != nil {
		user = &models.User{
			ID:         uuid.New(),
			Email:      oauthUser.Email,
			Username:   s.generateUsername(oauthUser.Name, oauthUser.Email),
			Avatar:     oauthUser.Avatar,
			Provider:   req.Provider,
			ProviderID: oauthUser.ID,
			IsActive:   true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}

		// Automatically create wallet for new user
		if s.walletService != nil {
			_, walletErr := s.walletService.CreateWallet(user.ID)
			if walletErr != nil {
				// Log the error but don't fail the auth process
				// In production, you might want to use a proper logger
				// log.Printf("Failed to create wallet for user %s: %v", user.ID, walletErr)
			}
		}
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	return s.userRepo.FindByID(ctx, userID)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid refresh token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *AuthService) generateAuthResponse(ctx context.Context, user *models.User) (*models.AuthResponse, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// Helper methods for password handling and username generation
func (s *AuthService) verifyPassword(password, hashedPassword string) bool {
	return utils.VerifyPassword(password, hashedPassword)
}

func (s *AuthService) hashPassword(password string) (string, error) {
	return utils.HashPassword(password)
}

func (s *AuthService) generateUsername(name, email string) string {
	if name != "" {
		return name
	}
	// Extract username from email (part before @)
	for i, char := range email {
		if char == '@' {
			return email[:i]
		}
	}
	return "user"
}
