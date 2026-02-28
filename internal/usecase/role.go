package usecase

import (
	"context"
	"log/slog"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"
)

// CreateRoleInput holds the fields required to create a new role.
type CreateRoleInput struct {
	Name        string
	Description string
	Permissions []domain.Permission
	CreatedBy   string
}

// UpdateRoleInput holds the fields required to update an existing role.
type UpdateRoleInput struct {
	ID          string
	Name        string
	Description string
	Permissions []domain.Permission
	UpdatedBy   string
}

// Role defines role use case operations
type Role interface {
	Create(ctx context.Context, in CreateRoleInput) (*domain.Role, error)
	Get(ctx context.Context, id string) (*domain.Role, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Role, error)
	Update(ctx context.Context, in UpdateRoleInput) (*domain.Role, error)
	Delete(ctx context.Context, id string) error
}

type role struct {
	repo repository.RoleRepository
}

// NewRole creates a new role use case
func NewRole(repo repository.RoleRepository) Role {
	return &role{
		repo: repo,
	}
}

// Create creates a new role
func (ref *role) Create(ctx context.Context, in CreateRoleInput) (*domain.Role, error) {
	if in.Name == "" {
		slog.Error("role creation failed: name is required")
		return nil, domain.ErrTitleRequired
	}

	r := domain.NewRole(in.Name, in.Description, in.Permissions, in.CreatedBy)

	if err := ref.repo.Create(ctx, r); err != nil {
		slog.Error("failed to persist role", "error", err, "name", in.Name)
		return nil, domain.ErrInternalServerError
	}

	return r, nil
}

// Get retrieves a role by ID
func (ref *role) Get(ctx context.Context, id string) (*domain.Role, error) {
	r, err := ref.repo.Get(ctx, id)
	if err != nil {
		slog.Error("failed to get role", "error", err, "roleID", id)
		return nil, domain.ErrNotFound
	}
	return r, nil
}

// List retrieves a list of roles
func (ref *role) List(ctx context.Context, limit, offset int) ([]*domain.Role, error) {
	if limit <= 0 {
		limit = 20
	}
	roles, err := ref.repo.List(ctx, limit, offset)
	if err != nil {
		slog.Error("failed to list roles", "error", err, "limit", limit, "offset", offset)
		return nil, domain.ErrInternalServerError
	}
	return roles, nil
}

// Update updates an existing role
func (ref *role) Update(ctx context.Context, in UpdateRoleInput) (*domain.Role, error) {
	r, err := ref.repo.Get(ctx, in.ID)
	if err != nil {
		slog.Error("failed to get role for update", "error", err, "roleID", in.ID)
		return nil, domain.ErrNotFound
	}

	r.Name = in.Name
	r.Description = in.Description
	r.Permissions = in.Permissions
	r.UpdateTimestamp(in.UpdatedBy)

	if err := ref.repo.Update(ctx, r); err != nil {
		slog.Error("failed to persist role update", "error", err, "roleID", in.ID)
		return nil, domain.ErrInternalServerError
	}

	return r, nil
}

// Delete deletes a role
func (ref *role) Delete(ctx context.Context, id string) error {
	if err := ref.repo.Delete(ctx, id); err != nil {
		slog.Error("failed to delete role", "error", err, "roleID", id)
		return domain.ErrInternalServerError
	}
	return nil
}
