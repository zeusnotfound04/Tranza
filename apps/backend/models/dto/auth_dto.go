package dto

// User Registration Request
type RegisterRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required" binding:"required"`
	Email       string `json:"email" validate:"email"`
	FullName    string `json:"full_name" validate:"required,min=2,max=100" binding:"required"`
	Password    string `json:"password" validate:"required,min=6" binding:"required"`
}

// User Login Request
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required" binding:"required"`
	Password    string `json:"password" validate:"required" binding:"required"`
}

// OTP Verification Request
type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required" binding:"required"`
	OTP         string `json:"otp" validate:"required,len=6" binding:"required"`
}

// Send OTP Request
type SendOTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required" binding:"required"`
}

// Login Response
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	TokenType    string       `json:"token_type"`
	User         UserResponse `json:"user"`
}

// User Response
type UserResponse struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	KYCStatus   string `json:"kyc_status"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}