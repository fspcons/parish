package handler

import (
	"encoding/json"
	"net/http"

	"github.com/parish/cmd/rest/middleware"
	"github.com/parish/internal/usecase"
)

const authCookieName = "auth_token"

const resetPasswordSuccessMessage = "Foi enviado um e-mail com instruções para o e-mail informado"

// AuthHandler handles authentication requests
type AuthHandler struct {
	authUseCase  usecase.Auth
	cookieSecure bool
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUseCase usecase.Auth, cookieSecure bool) *AuthHandler {
	return &AuthHandler{
		authUseCase:  authUseCase,
		cookieSecure: cookieSecure,
	}
}

// Register handles user registration
func (ref *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := ref.authUseCase.Register(r.Context(), usecase.RegisterInput{
		Email:     req.Email,
		Name:      req.Name,
		Password:  req.Password,
		RoleIDs:   req.RoleIDs,
		CreatedBy: "system",
	})
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusCreated, user.ToResponse(), "User registered successfully")
}

// Login handles user login
func (ref *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, token, err := ref.authUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	setAuthCookie(w, token, ref.cookieSecure)
	RespondSuccess(w, http.StatusOK, LoginResponse{User: user.ToResponse()}, "Login successful")
}

// Logout clears the authentication cookie
func (ref *AuthHandler) Logout(w http.ResponseWriter, _ *http.Request) {
	clearAuthCookie(w, ref.cookieSecure)
	RespondSuccess(w, http.StatusOK, nil, "Logged out successfully")
}

// ResetPassword triggers a password reset e-mail when the account exists. HTTP response is always the same.
func (ref *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if err := ref.authUseCase.RequestPasswordReset(r.Context(), req.Email); err != nil {
		HandleDomainError(w, err)
		return
	}
	RespondSuccess(w, http.StatusOK, nil, resetPasswordSuccessMessage)
}

// ChangePassword lets the authenticated user set a new password (clears forced change flag).
func (ref *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := ref.authUseCase.ChangePassword(r.Context(), user.ID, req.CurrentPassword, req.NewPassword); err != nil {
		HandleDomainError(w, err)
		return
	}
	RespondSuccess(w, http.StatusOK, nil, "Password updated successfully")
}

func setAuthCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   86400, // 24h
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

func clearAuthCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}
