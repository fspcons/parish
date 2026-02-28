package handler

import (
	"encoding/json"
	"net/http"

	"github.com/parish/internal/usecase"
)

const authCookieName = "auth_token"

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

	RespondSuccess(w, http.StatusCreated, user, "User registered successfully")
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
	RespondSuccess(w, http.StatusOK, LoginResponse{User: user}, "Login successful")
}

// Logout clears the authentication cookie
func (ref *AuthHandler) Logout(w http.ResponseWriter, _ *http.Request) {
	clearAuthCookie(w, ref.cookieSecure)
	RespondSuccess(w, http.StatusOK, nil, "Logged out successfully")
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
