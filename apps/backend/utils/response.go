package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StandardResponse represents the standard API response format
type StandardResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Code      string      `json:"code,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// PaginatedResponse represents paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Timestamp  string      `json:"timestamp"`
}

// Pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	c.JSON(statusCode, StandardResponse{
		Success:   false,
		Message:   message,
		Error:     errorMsg,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// ErrorResponseWithCode sends an error response with custom error code
func ErrorResponseWithCode(c *gin.Context, statusCode int, message string, errorCode string, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	c.JSON(statusCode, StandardResponse{
		Success:   false,
		Message:   message,
		Error:     errorMsg,
		Code:      errorCode,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// PaginatedSuccessResponse sends a paginated success response
func PaginatedSuccessResponse(c *gin.Context, message string, data interface{}, page, limit int, total int64) {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	
	pagination := Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})
}

// BadRequestResponse sends a 400 Bad Request response
func BadRequestResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusBadRequest, message, err)
}

// UnauthorizedResponse sends a 401 Unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

// ForbiddenResponse sends a 403 Forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, nil)
}

// NotFoundResponse sends a 404 Not Found response
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

// InternalServerErrorResponse sends a 500 Internal Server Error response
func InternalServerErrorResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusInternalServerError, message, err)
}

// ValidationErrorResponse sends a 422 Unprocessable Entity response for validation errors
func ValidationErrorResponse(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, StandardResponse{
		Success:   false,
		Message:   "Validation failed",
		Data:      errors,
		Code:      "VALIDATION_ERROR",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
