package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/parish/internal/cache"
	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"
)

const tokenTTL = 24 * time.Hour

// RegisterInput holds the fields required to register a new user.
type RegisterInput struct {
	Email     string
	Name      string
	Password  string
	RoleIDs   []string
	CreatedBy string
}

// Auth defines authentication use case operations
type Auth interface {
	Register(ctx context.Context, in RegisterInput) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
	GetUserPermissions(ctx context.Context, userID string) ([]domain.Permission, error)
	CheckPermission(ctx context.Context, userID, resource string, write bool) (bool, error)
	AssignRoles(ctx context.Context, userID string, roleIDs []string, updatedBy string) error
}

type auth struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	cache    cache.Cache
}

// NewAuth creates a new auth use case
func NewAuth(userRepo repository.UserRepository, roleRepo repository.RoleRepository, c cache.Cache) Auth {
	return &auth{
		userRepo: userRepo,
		roleRepo: roleRepo,
		cache:    c,
	}
}

// Register creates a new user
func (ref *auth) Register(ctx context.Context, in RegisterInput) (*domain.User, error) {
	if in.Email == "" {
		slog.Error("registration failed: email is required")
		return nil, domain.ErrEmailRequired
	}
	if in.Password == "" {
		slog.Error("registration failed: password is required")
		return nil, domain.ErrPasswordRequired
	}

	existingUser, _ := ref.userRepo.GetByEmail(ctx, in.Email)
	if existingUser != nil {
		slog.Error("registration failed: user already exists", "email", in.Email)
		return nil, domain.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, domain.ErrInternalServerError
	}

	user := domain.NewUser(in.Email, in.Name, string(hashedPassword), in.RoleIDs, in.CreatedBy)

	if err := ref.userRepo.Create(ctx, user); err != nil {
		slog.Error("failed to persist user", "error", err, "email", in.Email)
		return nil, domain.ErrInternalServerError
	}

	return user, nil
}

// Login authenticates a user and returns a token
func (ref *auth) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := ref.userRepo.GetByEmail(ctx, email)
	if err != nil {
		slog.Error("login failed: user not found", "email", email)
		return nil, "", domain.ErrInvalidCredentials
	}

	if !user.IsActive() {
		slog.Error("login failed: user is inactive", "email", email)
		return nil, "", domain.ErrUserInactive
	}

	if err := user.CheckPassword(password); err != nil {
		slog.Error("login failed: invalid password", "email", email)
		return nil, "", domain.ErrInvalidCredentials
	}

	// Generate token
	token, err := ref.generateToken(ctx, user.ID)
	if err != nil {
		slog.Error("failed to generate token", "error", err, "userID", user.ID)
		return nil, "", domain.ErrInternalServerError
	}

	return user, token, nil
}

// ValidateToken validates a token and returns the associated user
func (ref *auth) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	userID, err := ref.cache.Get(ctx, "token:"+token)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			slog.Error("token validation failed: token not found or expired")
			return nil, domain.ErrInvalidCredentials
		}
		slog.Error("token validation failed: cache error", "error", err)
		return nil, domain.ErrInternalServerError
	}

	user, err := ref.userRepo.Get(ctx, userID)
	if err != nil {
		slog.Error("token validation failed: user not found", "error", err, "userID", userID)
		return nil, domain.ErrNotFound
	}

	if !user.IsActive() {
		slog.Error("token validation failed: user is inactive", "userID", userID)
		return nil, domain.ErrUserInactive
	}

	return user, nil
}

// GetUserPermissions retrieves all permissions for a user
func (ref *auth) GetUserPermissions(ctx context.Context, userID string) ([]domain.Permission, error) {
	user, err := ref.userRepo.Get(ctx, userID)
	if err != nil {
		slog.Error("failed to get user for permissions", "error", err, "userID", userID)
		return nil, domain.ErrNotFound
	}

	roles, err := ref.roleRepo.GetMultiple(ctx, user.RoleIDs)
	if err != nil {
		slog.Error("failed to get roles for user", "error", err, "userID", userID)
		return nil, domain.ErrInternalServerError
	}

	// Merge permissions from all roles
	permMap := make(map[string]*domain.Permission)
	for _, role := range roles {
		if role == nil {
			continue
		}
		for _, perm := range role.Permissions {
			existing, exists := permMap[perm.Resource]
			if !exists {
				permCopy := perm
				permMap[perm.Resource] = &permCopy
			} else {
				// Merge permissions (if any role grants read/write, grant it)
				existing.Read = existing.Read || perm.Read
				existing.Write = existing.Write || perm.Write
			}
		}
	}

	// Convert map to slice
	permissions := make([]domain.Permission, 0, len(permMap))
	for _, perm := range permMap {
		permissions = append(permissions, *perm)
	}

	return permissions, nil
}

// CheckPermission checks if a user has permission for a resource
func (ref *auth) CheckPermission(ctx context.Context, userID, resource string, write bool) (bool, error) {
	permissions, err := ref.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if perm.Resource == resource {
			if write {
				return perm.Write, nil
			}
			return perm.Read, nil
		}
	}

	return false, nil
}

// AssignRoles assigns the given role IDs to a user
func (ref *auth) AssignRoles(ctx context.Context, userID string, roleIDs []string, updatedBy string) error {
	user, err := ref.userRepo.Get(ctx, userID)
	if err != nil {
		slog.Error("failed to get user for role assignment", "error", err, "userID", userID)
		return domain.ErrNotFound
	}

	// Validate that all role IDs exist
	roles, err := ref.roleRepo.GetMultiple(ctx, roleIDs)
	if err != nil {
		slog.Error("failed to validate role IDs", "error", err, "roleIDs", roleIDs)
		return domain.ErrInternalServerError
	}
	for i, r := range roles {
		if r == nil {
			slog.Error("role not found for assignment", "roleID", roleIDs[i])
			return domain.ErrNotFound
		}
	}

	user.RoleIDs = roleIDs
	user.UpdateTimestamp(updatedBy)

	if err := ref.userRepo.Update(ctx, user); err != nil {
		slog.Error("failed to persist user role assignment", "error", err, "userID", userID)
		return domain.ErrInternalServerError
	}

	slog.Info("roles assigned to user", "userID", userID, "roleIDs", roleIDs)
	return nil
}

// generateToken generates a random token and stores it in the cache.
func (ref *auth) generateToken(ctx context.Context, userID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(b)

	if err := ref.cache.Set(ctx, "token:"+token, userID, tokenTTL); err != nil {
		return "", err
	}

	return token, nil
}
