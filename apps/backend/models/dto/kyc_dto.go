package dto

// KYC Upload Request
type KYCUploadRequest struct {
	DocumentType   string `json:"document_type" validate:"required,oneof=aadhar pan driving_license passport" binding:"required"`
	DocumentNumber string `json:"document_number" validate:"required,min=5,max=50" binding:"required"`
	DocumentFront  string `json:"document_front" validate:"required"` // Base64 encoded image
	DocumentBack   string `json:"document_back,omitempty"`            // Base64 encoded image (for driving license)
}

// KYC Status Response
type KYCStatusResponse struct {
	Status         string `json:"status"`
	DocumentType   string `json:"document_type,omitempty"`
	DocumentNumber string `json:"document_number,omitempty"`
	SubmittedAt    string `json:"submitted_at,omitempty"`
	VerifiedAt     string `json:"verified_at,omitempty"`
	RejectedAt     string `json:"rejected_at,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	Message        string `json:"message"`
}
