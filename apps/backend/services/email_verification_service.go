package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// EmailVerificationService handles email verification business logic
type EmailVerificationService struct {
	emailVerificationRepo *repositories.EmailVerificationRepository
	userRepo              repositories.UserRepository
	emailService          *EmailService
}

// NewEmailVerificationService creates a new email verification service
func NewEmailVerificationService(
	emailVerificationRepo *repositories.EmailVerificationRepository,
	userRepo repositories.UserRepository,
	emailService *EmailService,
) *EmailVerificationService {
	return &EmailVerificationService{
		emailVerificationRepo: emailVerificationRepo,
		userRepo:              userRepo,
		emailService:          emailService,
	}
}

// Constants for verification
const (
	MaxVerificationAttempts = 5
	VerificationExpiration  = 15 * time.Minute
	ResendCooldown         = 2 * time.Minute
)

var (
	ErrEmailAlreadyExists        = errors.New("email already registered")
	ErrUsernameAlreadyExists     = errors.New("username already taken")
	ErrVerificationNotFound      = errors.New("verification record not found")
	ErrVerificationExpired       = errors.New("verification code expired")
	ErrInvalidVerificationCode   = errors.New("invalid verification code")
	ErrTooManyAttempts          = errors.New("too many verification attempts")
	ErrResendCooldown           = errors.New("please wait before requesting another code")
)

// InitiateEmailVerification starts the email verification process
func (evs *EmailVerificationService) InitiateEmailVerification(ctx context.Context, req models.PreRegistrationRequest) (*models.PreRegistrationResponse, error) {
	// Check if email already exists
	emailExists, err := evs.emailVerificationRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		return nil, ErrEmailAlreadyExists
	}

	// Check if username already exists
	usernameExists, err := evs.emailVerificationRepo.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if usernameExists {
		return nil, ErrUsernameAlreadyExists
	}

	// Check if there's an existing verification record
	existingVerification, err := evs.emailVerificationRepo.GetVerificationByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing verification: %w", err)
	}

	// If verification exists and is recent, check cooldown
	if existingVerification != nil {
		if time.Since(existingVerification.UpdatedAt) < ResendCooldown {
			return nil, ErrResendCooldown
		}
		// Delete old verification record
		if err := evs.emailVerificationRepo.DeleteVerification(ctx, req.Email); err != nil {
			return nil, fmt.Errorf("failed to delete old verification: %w", err)
		}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate verification code
	code, err := evs.emailService.GenerateVerificationCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Hash the verification code for storage
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash verification code: %w", err)
	}

	// Create verification record
	verification := &models.EmailVerification{
		Email:      req.Email,
		Username:   req.Username,
		Password:   string(hashedPassword),
		Code:       string(hashedCode),
		ExpiresAt:  time.Now().Add(VerificationExpiration),
		Attempts:   0,
		IsVerified: false,
	}

	if err := evs.emailVerificationRepo.CreateVerification(ctx, verification); err != nil {
		return nil, fmt.Errorf("failed to create verification record: %w", err)
	}

	// Send verification email
	if err := evs.emailService.SendVerificationEmail(req.Email, req.Username, code); err != nil {
		return nil, fmt.Errorf("failed to send verification email: %w", err)
	}

	return &models.PreRegistrationResponse{
		Message:   "Verification code sent to your email",
		Email:     req.Email,
		ExpiresAt: verification.ExpiresAt,
	}, nil
}

// VerifyEmailCode verifies the email verification code and creates the user
func (evs *EmailVerificationService) VerifyEmailCode(ctx context.Context, req models.EmailVerificationRequest) (*models.EmailVerificationResponse, error) {
	// Get verification record
	verification, err := evs.emailVerificationRepo.GetVerificationByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVerificationNotFound
		}
		return nil, fmt.Errorf("failed to get verification record: %w", err)
	}

	// Check if verification has expired
	if time.Now().After(verification.ExpiresAt) {
		// Clean up expired record
		evs.emailVerificationRepo.DeleteVerification(ctx, req.Email)
		return nil, ErrVerificationExpired
	}

	// Check attempts limit
	if verification.Attempts >= MaxVerificationAttempts {
		return nil, ErrTooManyAttempts
	}

	// Verify the code
	if err := bcrypt.CompareHashAndPassword([]byte(verification.Code), []byte(req.Code)); err != nil {
		// Increment attempts
		evs.emailVerificationRepo.IncrementAttempts(ctx, req.Email)
		return nil, ErrInvalidVerificationCode
	}

	// Create the user account
	user := &models.User{
		Email:     verification.Email,
		Username:  verification.Username,
		Password:  verification.Password, // Already hashed
		Provider:  "local",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := evs.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Mark verification as completed and delete the record
	if err := evs.emailVerificationRepo.DeleteVerification(ctx, req.Email); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to delete verification record: %v\n", err)
	}

	// Send welcome email
	if err := evs.emailService.SendWelcomeEmail(user.Email, user.Username); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to send welcome email: %v\n", err)
	}

	return &models.EmailVerificationResponse{
		Message: "Email verified successfully! Account created.",
		User:    user,
	}, nil
}

// ResendVerificationCode resends the verification code
func (evs *EmailVerificationService) ResendVerificationCode(ctx context.Context, req models.ResendVerificationRequest) (*models.PreRegistrationResponse, error) {
	// Get verification record
	verification, err := evs.emailVerificationRepo.GetVerificationByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVerificationNotFound
		}
		return nil, fmt.Errorf("failed to get verification record: %w", err)
	}

	// Check cooldown period
	if time.Since(verification.UpdatedAt) < ResendCooldown {
		return nil, ErrResendCooldown
	}

	// Check if verification has expired
	if time.Now().After(verification.ExpiresAt) {
		// Delete expired record
		evs.emailVerificationRepo.DeleteVerification(ctx, req.Email)
		return nil, ErrVerificationExpired
	}

	// Generate new verification code
	code, err := evs.emailService.GenerateVerificationCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Hash the new verification code
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash verification code: %w", err)
	}

	// Update verification record
	verification.Code = string(hashedCode)
	verification.ExpiresAt = time.Now().Add(VerificationExpiration)
	verification.Attempts = 0 // Reset attempts
	verification.UpdatedAt = time.Now()

	if err := evs.emailVerificationRepo.UpdateVerification(ctx, verification); err != nil {
		return nil, fmt.Errorf("failed to update verification record: %w", err)
	}

	// Send verification email
	if err := evs.emailService.SendVerificationEmail(req.Email, verification.Username, code); err != nil {
		return nil, fmt.Errorf("failed to send verification email: %w", err)
	}

	return &models.PreRegistrationResponse{
		Message:   "New verification code sent to your email",
		Email:     req.Email,
		ExpiresAt: verification.ExpiresAt,
	}, nil
}

// CleanupExpiredVerifications removes expired verification records
func (evs *EmailVerificationService) CleanupExpiredVerifications(ctx context.Context) error {
	return evs.emailVerificationRepo.DeleteExpiredVerifications(ctx)
}