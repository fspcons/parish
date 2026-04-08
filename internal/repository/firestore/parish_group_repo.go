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

// ParishGroupRepository implements repository.ParishGroupRepository.
type ParishGroupRepository struct {
	store *Store
}

// NewParishGroupRepository creates a ParishGroupRepository.
func NewParishGroupRepository(store *Store) *ParishGroupRepository {
	return &ParishGroupRepository{store: store}
}

func (r *ParishGroupRepository) col() *gcfs.CollectionRef {
	return r.store.Collection(colParishGroups)
}

// Create creates a parish group.
func (r *ParishGroupRepository) Create(ctx context.Context, group *domain.ParishGroup) error {
	_, err := r.col().Doc(group.ID).Create(ctx, group)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return fmt.Errorf("parish group already exists: %s", group.ID)
		}
		return fmt.Errorf("failed to create parish group: %w", err)
	}
	return nil
}

// Get returns a parish group by ID.
func (r *ParishGroupRepository) Get(ctx context.Context, id string) (*domain.ParishGroup, error) {
	snap, err := r.col().Doc(id).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("parish group not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get parish group: %w", err)
	}
	var g domain.ParishGroup
	if err := snap.DataTo(&g); err != nil {
		return nil, fmt.Errorf("failed to decode parish group: %w", err)
	}
	g.ID = id
	return &g, nil
}

// List lists parish groups by createdAt descending.
func (r *ParishGroupRepository) List(ctx context.Context, limit, offset int) ([]*domain.ParishGroup, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	q := r.col().OrderBy("createdAt", gcfs.Desc).Limit(limit).Offset(offset)
	iter := q.Documents(ctx)
	

	return scanDocuments[domain.ParishGroup](iter)
}

// ListByDateRange lists groups created in [start, end].
func (r *ParishGroupRepository) ListByDateRange(ctx context.Context, start, end time.Time) ([]*domain.ParishGroup, error) {
	q := r.col().
		Where("createdAt", ">=", start).
		Where("createdAt", "<=", end).
		OrderBy("createdAt", gcfs.Asc)
	iter := q.Documents(ctx)
	

	return scanDocuments[domain.ParishGroup](iter)
}

// Update saves a parish group.
func (r *ParishGroupRepository) Update(ctx context.Context, group *domain.ParishGroup) error {
	_, err := r.col().Doc(group.ID).Set(ctx, group)
	if err != nil {
		return fmt.Errorf("failed to update parish group: %w", err)
	}
	return nil
}

// Delete removes a parish group.
func (r *ParishGroupRepository) Delete(ctx context.Context, id string) error {
	_, err := r.col().Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete parish group: %w", err)
	}
	return nil
}
