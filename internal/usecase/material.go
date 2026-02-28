package usecase

import (
	"context"
	"log/slog"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"
)

// CreateMaterialInput holds the fields required to create a new material.
type CreateMaterialInput struct {
	Title       string
	Type        string
	Description string
	URL         string
	Label       string
	CreatedBy   string
}

// UpdateMaterialInput holds the fields required to update an existing material.
type UpdateMaterialInput struct {
	ID          string
	Title       string
	Type        string
	Description string
	URL         string
	Label       string
	UpdatedBy   string
}

// Material defines material use case operations
type Material interface {
	Create(ctx context.Context, in CreateMaterialInput) (*domain.Material, error)
	Get(ctx context.Context, id string) (*domain.Material, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Material, error)
	ListByType(ctx context.Context, materialType string, limit, offset int) ([]*domain.Material, error)
	ListByLabel(ctx context.Context, label string, limit, offset int) ([]*domain.Material, error)
	Update(ctx context.Context, in UpdateMaterialInput) (*domain.Material, error)
	Delete(ctx context.Context, id string) error
}

type material struct {
	repo repository.MaterialRepository
}

// NewMaterial creates a new material use case
func NewMaterial(repo repository.MaterialRepository) Material {
	return &material{
		repo: repo,
	}
}

// Create creates a new material
func (ref *material) Create(ctx context.Context, in CreateMaterialInput) (*domain.Material, error) {
	mat, err := domain.NewMaterial(in.Title, in.Type, in.Description, in.URL, in.Label, in.CreatedBy)
	if err != nil {
		slog.Error("failed to create material entity", "error", err, "title", in.Title, "type", in.Type)
		return nil, err
	}

	if err := ref.repo.Create(ctx, mat); err != nil {
		slog.Error("failed to persist material", "error", err, "materialID", mat.ID)
		return nil, domain.ErrInternalServerError
	}

	return mat, nil
}

// Get retrieves a material by ID
func (ref *material) Get(ctx context.Context, id string) (*domain.Material, error) {
	mat, err := ref.repo.Get(ctx, id)
	if err != nil {
		slog.Error("failed to get material", "error", err, "materialID", id)
		return nil, domain.ErrNotFound
	}
	return mat, nil
}

// List retrieves a list of materials
func (ref *material) List(ctx context.Context, limit, offset int) ([]*domain.Material, error) {
	if limit <= 0 {
		limit = 20
	}
	materials, err := ref.repo.List(ctx, limit, offset)
	if err != nil {
		slog.Error("failed to list materials", "error", err, "limit", limit, "offset", offset)
		return nil, domain.ErrInternalServerError
	}
	return materials, nil
}

// ListByType retrieves materials by type
func (ref *material) ListByType(ctx context.Context, materialType string, limit, offset int) ([]*domain.Material, error) {
	if limit <= 0 {
		limit = 20
	}
	materials, err := ref.repo.ListByType(ctx, materialType, limit, offset)
	if err != nil {
		slog.Error("failed to list materials by type", "error", err, "type", materialType)
		return nil, domain.ErrInternalServerError
	}
	return materials, nil
}

// ListByLabel retrieves materials by label
func (ref *material) ListByLabel(ctx context.Context, label string, limit, offset int) ([]*domain.Material, error) {
	if limit <= 0 {
		limit = 20
	}
	materials, err := ref.repo.ListByLabel(ctx, label, limit, offset)
	if err != nil {
		slog.Error("failed to list materials by label", "error", err, "label", label)
		return nil, domain.ErrInternalServerError
	}
	return materials, nil
}

// Update updates an existing material
func (ref *material) Update(ctx context.Context, in UpdateMaterialInput) (*domain.Material, error) {
	mat, err := ref.repo.Get(ctx, in.ID)
	if err != nil {
		slog.Error("failed to get material for update", "error", err, "materialID", in.ID)
		return nil, domain.ErrNotFound
	}

	if err := mat.Update(in.Title, in.Type, in.Description, in.URL, in.Label, in.UpdatedBy); err != nil {
		slog.Error("failed to update material entity", "error", err, "materialID", in.ID)
		return nil, err
	}

	if err := ref.repo.Update(ctx, mat); err != nil {
		slog.Error("failed to persist material update", "error", err, "materialID", in.ID)
		return nil, domain.ErrInternalServerError
	}

	return mat, nil
}

// Delete deletes a material
func (ref *material) Delete(ctx context.Context, id string) error {
	if err := ref.repo.Delete(ctx, id); err != nil {
		slog.Error("failed to delete material", "error", err, "materialID", id)
		return domain.ErrInternalServerError
	}
	return nil
}
