package usecase

import (
	"context"
	"log/slog"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"
)

// CreateParishGroupInput holds the fields required to create a new parish group.
type CreateParishGroupInput struct {
	Title       string
	Description string
	Manager     string
	Active      bool
	CreatedBy   string
}

// UpdateParishGroupInput holds the fields required to update a parish group.
type UpdateParishGroupInput struct {
	ID          string
	Title       string
	Description string
	Manager     string
	Active      bool
	UpdatedBy   string
}

// ParishGroup defines parish group use case operations
type ParishGroup interface {
	Create(ctx context.Context, in CreateParishGroupInput) (*domain.ParishGroup, error)
	Get(ctx context.Context, id string) (*domain.ParishGroup, error)
	List(ctx context.Context, limit, offset int) ([]*domain.ParishGroup, error)
	Update(ctx context.Context, in UpdateParishGroupInput) (*domain.ParishGroup, error)
	Delete(ctx context.Context, id string) error
}

type parishGroup struct {
	repo repository.ParishGroupRepository
}

// NewParishGroup creates a new parish group use case
func NewParishGroup(repo repository.ParishGroupRepository) ParishGroup {
	return &parishGroup{
		repo: repo,
	}
}

// Create creates a new parish group
func (ref *parishGroup) Create(ctx context.Context, in CreateParishGroupInput) (*domain.ParishGroup, error) {
	group, err := domain.NewParishGroup(in.Title, in.Description, in.Manager, in.Active, in.CreatedBy)
	if err != nil {
		slog.Error("failed to create parish group entity", "error", err, "title", in.Title)
		return nil, err
	}

	if err := ref.repo.Create(ctx, group); err != nil {
		slog.Error("failed to persist parish group", "error", err, "groupID", group.ID)
		return nil, domain.ErrInternalServerError
	}

	return group, nil
}

// Get retrieves a parish group by ID
func (ref *parishGroup) Get(ctx context.Context, id string) (*domain.ParishGroup, error) {
	group, err := ref.repo.Get(ctx, id)
	if err != nil {
		slog.Error("failed to get parish group", "error", err, "groupID", id)
		return nil, domain.ErrNotFound
	}
	return group, nil
}

// List retrieves a list of parish groups
func (ref *parishGroup) List(ctx context.Context, limit, offset int) ([]*domain.ParishGroup, error) {
	if limit <= 0 {
		limit = 20
	}
	groups, err := ref.repo.List(ctx, limit, offset)
	if err != nil {
		slog.Error("failed to list parish groups", "error", err, "limit", limit, "offset", offset)
		return nil, domain.ErrInternalServerError
	}
	return groups, nil
}

// Update updates an existing parish group
func (ref *parishGroup) Update(ctx context.Context, in UpdateParishGroupInput) (*domain.ParishGroup, error) {
	group, err := ref.repo.Get(ctx, in.ID)
	if err != nil {
		slog.Error("failed to get parish group for update", "error", err, "groupID", in.ID)
		return nil, domain.ErrNotFound
	}

	if err := group.Update(in.Title, in.Description, in.Manager, in.Active, in.UpdatedBy); err != nil {
		slog.Error("failed to update parish group entity", "error", err, "groupID", in.ID)
		return nil, err
	}

	if err := ref.repo.Update(ctx, group); err != nil {
		slog.Error("failed to persist parish group update", "error", err, "groupID", in.ID)
		return nil, domain.ErrInternalServerError
	}

	return group, nil
}

// Delete deletes a parish group
func (ref *parishGroup) Delete(ctx context.Context, id string) error {
	if err := ref.repo.Delete(ctx, id); err != nil {
		slog.Error("failed to delete parish group", "error", err, "groupID", id)
		return domain.ErrInternalServerError
	}
	return nil
}
