package firestore

import (
	"context"
	"fmt"
	"time"

	gcfs "cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/parish/internal/domain"
)

// MaterialRepository implements repository.MaterialRepository.
type MaterialRepository struct {
	store *Store
}

// NewMaterialRepository creates a MaterialRepository.
func NewMaterialRepository(store *Store) *MaterialRepository {
	return &MaterialRepository{store: store}
}

func (r *MaterialRepository) col() *gcfs.CollectionRef {
	return r.store.Collection(colMaterials)
}

// Create creates a material.
func (r *MaterialRepository) Create(ctx context.Context, material *domain.Material) error {
	_, err := r.col().Doc(material.ID).Create(ctx, material)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return fmt.Errorf("material already exists: %s", material.ID)
		}
		return fmt.Errorf("failed to create material: %w", err)
	}
	return nil
}

// Get returns a material by ID.
func (r *MaterialRepository) Get(ctx context.Context, id string) (*domain.Material, error) {
	snap, err := r.col().Doc(id).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("material not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get material: %w", err)
	}
	var m domain.Material
	if err := snap.DataTo(&m); err != nil {
		return nil, fmt.Errorf("failed to decode material: %w", err)
	}
	m.ID = id
	return &m, nil
}

// List lists materials by createdAt descending.
func (r *MaterialRepository) List(ctx context.Context, limit, offset int) ([]*domain.Material, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	q := r.col().OrderBy("createdAt", gcfs.Desc).Limit(limit).Offset(offset)
	iter := q.Documents(ctx)
	

	return scanDocuments[domain.Material](iter)
}

// ListByType lists materials filtered by type.
func (r *MaterialRepository) ListByType(ctx context.Context, materialType string, limit, offset int) ([]*domain.Material, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	q := r.col().
		Where("type", "==", materialType).
		OrderBy("createdAt", gcfs.Desc).
		Limit(limit).
		Offset(offset)
	iter := q.Documents(ctx)
	

	return scanDocuments[domain.Material](iter)
}

// ListByLabel lists materials filtered by label.
func (r *MaterialRepository) ListByLabel(ctx context.Context, label string, limit, offset int) ([]*domain.Material, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	q := r.col().
		Where("label", "==", label).
		OrderBy("createdAt", gcfs.Desc).
		Limit(limit).
		Offset(offset)
	iter := q.Documents(ctx)
	

	return scanDocuments[domain.Material](iter)
}

// ListByDateRange lists materials created in [start, end].
func (r *MaterialRepository) ListByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Material, error) {
	q := r.col().
		Where("createdAt", ">=", start).
		Where("createdAt", "<=", end).
		OrderBy("createdAt", gcfs.Asc)
	iter := q.Documents(ctx)
	

	return scanDocuments[domain.Material](iter)
}

// Update saves a material.
func (r *MaterialRepository) Update(ctx context.Context, material *domain.Material) error {
	_, err := r.col().Doc(material.ID).Set(ctx, material)
	if err != nil {
		return fmt.Errorf("failed to update material: %w", err)
	}
	return nil
}

// Delete removes a material.
func (r *MaterialRepository) Delete(ctx context.Context, id string) error {
	_, err := r.col().Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete material: %w", err)
	}
	return nil
}
