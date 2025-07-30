package services

import (
	"context"
	"errors"

	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
	"github.com/zeusnotfound04/Tranza/utils"
)

type AuthService struct {
	userRepo     repositories.UserRepository
	jwtService   utils.JWTService
	oauthService OAuthService
}

func NewAuthService(userRepo repositories.UserRepository, jwtService utils.JWTService, oauthService OAuthService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		jwtService:   jwtService,
		oauthService: oauthService,
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
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Provider: "local",
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
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
			Email:      oauthUser.Email,
			Username:   s.generateUsername(oauthUser.Name, oauthUser.Email),
			Avatar:     oauthUser.Avatar,
			Provider:   req.Provider,
			ProviderID: oauthUser.ID,
			IsActive:   true,
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return s.userRepo.FindByID(ctx, uint(userID))
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(ctx, uint(userID))
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
