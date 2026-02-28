package datastore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/parish/internal/domain"
)

// ParishGroupRepository implements repository.ParishGroupRepository
type ParishGroupRepository struct {
	client *Client
}

// NewParishGroupRepository creates a new ParishGroupRepository
func NewParishGroupRepository(client *Client) *ParishGroupRepository {
	return &ParishGroupRepository{
		client: client,
	}
}

// Create creates a new parish group
func (ref *ParishGroupRepository) Create(ctx context.Context, group *domain.ParishGroup) error {
	key := datastore.NameKey(group.EntityKind(), group.ID, nil)

	_, err := ref.client.client.Put(ctx, key, group)
	if err != nil {
		return fmt.Errorf("failed to create parish group: %w", err)
	}

	return nil
}

// Get retrieves a parish group by ID
func (ref *ParishGroupRepository) Get(ctx context.Context, id string) (*domain.ParishGroup, error) {
	key := datastore.NameKey((&domain.ParishGroup{}).EntityKind(), id, nil)
	group := &domain.ParishGroup{}

	err := ref.client.client.Get(ctx, key, group)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, fmt.Errorf("parish group not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get parish group: %w", err)
	}

	group.ID = key.Name
	return group, nil
}

// List retrieves a list of parish groups
func (ref *ParishGroupRepository) List(ctx context.Context, limit, offset int) ([]*domain.ParishGroup, error) {
	query := datastore.NewQuery((&domain.ParishGroup{}).EntityKind()).
		Order("-createdAt").
		Limit(limit).
		Offset(offset)

	var groups []*domain.ParishGroup
	keys, err := ref.client.client.GetAll(ctx, query, &groups)
	if err != nil {
		return nil, fmt.Errorf("failed to list parish groups: %w", err)
	}

	for i, key := range keys {
		groups[i].ID = key.Name
	}

	return groups, nil
}

// ListByDateRange retrieves parish groups created within a date range
func (ref *ParishGroupRepository) ListByDateRange(ctx context.Context, start, end time.Time) ([]*domain.ParishGroup, error) {
	query := datastore.NewQuery((&domain.ParishGroup{}).EntityKind()).
		FilterField("createdAt", ">=", start).
		FilterField("createdAt", "<=", end).
		Order("createdAt")

	var groups []*domain.ParishGroup
	keys, err := ref.client.client.GetAll(ctx, query, &groups)
	if err != nil {
		return nil, fmt.Errorf("failed to list parish groups by date range: %w", err)
	}

	for i, key := range keys {
		groups[i].ID = key.Name
	}

	return groups, nil
}

// Update updates an existing parish group
func (ref *ParishGroupRepository) Update(ctx context.Context, group *domain.ParishGroup) error {
	key := datastore.NameKey(group.EntityKind(), group.ID, nil)

	_, err := ref.client.client.Put(ctx, key, group)
	if err != nil {
		return fmt.Errorf("failed to update parish group: %w", err)
	}

	return nil
}

// Delete deletes a parish group by ID
func (ref *ParishGroupRepository) Delete(ctx context.Context, id string) error {
	key := datastore.NameKey((&domain.ParishGroup{}).EntityKind(), id, nil)

	err := ref.client.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete parish group: %w", err)
	}

	return nil
}
