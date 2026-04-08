package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/parish/internal/domain"
)

type (
	// RegisterRequest represents a registration request
	RegisterRequest struct {
		Email    string   `json:"email"`
		Name     string   `json:"name"`
		Password string   `json:"password"`
		RoleIDs  []string `json:"roleIds"`
	}

	// LoginRequest represents a login request
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// LoginResponse represents a login response
	LoginResponse struct {
		User domain.UserResponse `json:"user"`
	}

	// ResetPasswordRequest requests a temporary password by e-mail (response is always the same).
	ResetPasswordRequest struct {
		Email string `json:"email"`
	}

	// ChangePasswordRequest is used by an authenticated user to set a new password.
	ChangePasswordRequest struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	// EventRequest represents an event request
	EventRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ImgURL      string `json:"imgUrl"`
		Date        string `json:"date"`
		Location    string `json:"location"`
		Origin      string `json:"origin"`
	}

	// MaterialRequest represents a material request
	MaterialRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Type        string `json:"type"`
		URL         string `json:"url"`
		Label       string `json:"label"`
	}

	// ParishGroupRequest represents a parish group request
	ParishGroupRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Manager     string `json:"manager"`
		Active      bool   `json:"active"`
	}

	// RoleRequest represents a role create/update request
	RoleRequest struct {
		Name        string              `json:"name"`
		Description string              `json:"description"`
		Permissions []domain.Permission `json:"permissions"`
	}

	// AssignRolesRequest represents a request to assign roles to a user
	AssignRolesRequest struct {
		RoleIDs []string `json:"roleIds"`
	}

	// ErrorResponse represents an error response
	ErrorResponse struct {
		Error string `json:"error"`
	}

	// SuccessResponse represents a success response
	SuccessResponse struct {
		Data    any    `json:"data,omitempty"`
		Message string `json:"message,omitempty"`
	}
)

// RespondJSON writes a JSON response
func RespondJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// RespondSuccess writes a success response
func RespondSuccess(w http.ResponseWriter, statusCode int, data any, message string) {
	RespondJSON(w, statusCode, SuccessResponse{Data: data, Message: message})
}

// HandleDomainError maps a domain error to the appropriate HTTP status code and writes a JSON error response.
func HandleDomainError(w http.ResponseWriter, err error) {
	errMsg := err.Error()
	var status int
	switch {
	case strings.HasPrefix(errMsg, "ERR_INVALID_ARGUMENT:"):
		status = http.StatusBadRequest
	case strings.HasPrefix(errMsg, "ERR_UNAUTHORIZED:"):
		status = http.StatusUnauthorized
	case strings.HasPrefix(errMsg, "ERR_FORBIDDEN:"):
		status = http.StatusForbidden
	case strings.HasPrefix(errMsg, "ERR_NOT_FOUND:"):
		status = http.StatusNotFound
	case strings.HasPrefix(errMsg, "ERR_CONFLICT:"):
		status = http.StatusConflict
	default:
		status = http.StatusInternalServerError
		errMsg = "internal error"
	}
	RespondJSON(w, status, ErrorResponse{Error: errMsg})
}
