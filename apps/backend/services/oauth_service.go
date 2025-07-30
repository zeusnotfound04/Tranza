package services

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/zeusnotfound04/Tranza/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type OAuthUser struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type OAuthService interface {
	GetAuthURL(provider, state, redirectURI string) (string, error)
	ExchangeCodeForUser(ctx context.Context, provider, code, redirectURI string) (*OAuthUser, error)
}

type oauthService struct {
	googleConfig *oauth2.Config
	githubConfig *oauth2.Config
}

func NewOAuthService(cfg *config.Config) OAuthService {
	return &oauthService{
		googleConfig: &oauth2.Config{
			ClientID:     cfg.OAuth.Google.ClientID,
			ClientSecret: cfg.OAuth.Google.ClientSecret,
			RedirectURL:  cfg.OAuth.Google.RedirectURL,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
		},
		githubConfig: &oauth2.Config{
			ClientID:     cfg.OAuth.GitHub.ClientID,
			ClientSecret: cfg.OAuth.GitHub.ClientSecret,
			RedirectURL:  cfg.OAuth.GitHub.RedirectURL,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

// NewOAuthServiceFromEnv creates a new OAuth service using environment variables
func NewOAuthServiceFromEnv() OAuthService {
	cfg := config.LoadOAuthConfig()
	return NewOAuthService(cfg)
}

func (s *oauthService) GetAuthURL(provider, state, redirectURI string) (string, error) {
	switch provider {
	case "google":
		return s.googleConfig.AuthCodeURL(state), nil
	case "github":
		return s.githubConfig.AuthCodeURL(state), nil
	default:
		return "", errors.New("unsupported provider")
	}
}

func (s *oauthService) ExchangeCodeForUser(ctx context.Context, provider, code, redirectURI string) (*OAuthUser, error) {
	switch provider {
	case "google":
		return s.handleGoogleAuth(ctx, code)
	case "github":
		return s.handleGitHubAuth(ctx, code)
	default:
		return nil, errors.New("unsupported provider")
	}
}

func (s *oauthService) handleGoogleAuth(ctx context.Context, code string) (*OAuthUser, error) {
	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := s.googleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &OAuthUser{
		ID:     googleUser.ID,
		Email:  googleUser.Email,
		Name:   googleUser.Name,
		Avatar: googleUser.Picture,
	}, nil
}

func (s *oauthService) handleGitHubAuth(ctx context.Context, code string) (*OAuthUser, error) {
	token, err := s.githubConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := s.githubConfig.Client(ctx, token)

	// Get user info
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID     int    `json:"id"`
		Login  string `json:"login"`
		Name   string `json:"name"`
		Avatar string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, err
	}

	// Get email (might be private)
	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return nil, err
	}
	defer emailResp.Body.Close()

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		return nil, err
	}

	var primaryEmail string
	for _, email := range emails {
		if email.Primary {
			primaryEmail = email.Email
			break
		}
	}

	return &OAuthUser{
		ID:     strconv.Itoa(githubUser.ID),
		Email:  primaryEmail,
		Name:   githubUser.Name,
		Avatar: githubUser.Avatar,
	}, nil
}
