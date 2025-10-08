package utils

import (
	"errors"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Error codes
const (
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeInvalidInput  = "INVALID_INPUT"
	ErrCodeInternal      = "INTERNAL_ERROR"
	ErrCodeDatabaseError = "DATABASE_ERROR"
	ErrCodeConflict      = "CONFLICT"
	ErrCodeUnauthorized  = "UNAUTHORIZED"
	ErrCodeForbidden     = "FORBIDDEN"
	ErrCodeNotAllowed    = "NOT_ALLOWED"
)

// AppError wraps errors with HTTP status codes and error codes
type AppError struct {
	Err        error
	Message    string
	Code       string
	StatusCode int
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(err error, message string, code string, statusCode int) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
	}
}

// Helper function to check if error is AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// APIResponse and ErrorInfo remain the same
type APIResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

func SuccessResponse(message string, data any) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string, err error) APIResponse {
	response := APIResponse{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = &ErrorInfo{
			Details: err.Error(),
		}
	}

	return response
}

func RespondWithError(ctx *gin.Context, err error, defaultMessage string) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		response := ErrorResponse(appErr.Message, appErr.Err)
		if response.Error != nil {
			response.Error.Code = appErr.Code
		}
		ctx.JSON(appErr.StatusCode, response)
		return
	}

	// fallback
	response := ErrorResponse(defaultMessage, err)
	if response.Error != nil {
		response.Error.Code = ErrCodeInternal
	}
	ctx.JSON(http.StatusInternalServerError, response)
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	default:
		return "invalid value"
	}
}

func RespondWithValidationError(ctx *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		details := make(map[string]string)
		for _, fe := range ve {
			details[fe.Field()] = getValidationMessage(fe)
		}

		// Convert details to string
		detailsStr := ""
		for field, msg := range details {
			if detailsStr != "" {
				detailsStr += "; "
			}
			detailsStr += field + ": " + msg
		}

		response := ErrorResponse("Validation failed", errors.New(detailsStr))
		if response.Error != nil {
			response.Error.Code = ErrCodeValidation
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func PaginatedResponse(message string, data interface{}, page, limit, total int) APIResponse {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return APIResponse{
		Success: true,
		Message: message,
		Data: map[string]interface{}{
			"items": data,
			"meta": PaginationMeta{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	}
}
