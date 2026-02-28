package datastore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/parish/internal/domain"
)

// MaterialRepository implements repository.MaterialRepository
type MaterialRepository struct {
	client *Client
}

// NewMaterialRepository creates a new MaterialRepository
func NewMaterialRepository(client *Client) *MaterialRepository {
	return &MaterialRepository{
		client: client,
	}
}

// Create creates a new material
func (ref *MaterialRepository) Create(ctx context.Context, material *domain.Material) error {
	key := datastore.NameKey(material.EntityKind(), material.ID, nil)

	_, err := ref.client.client.Put(ctx, key, material)
	if err != nil {
		return fmt.Errorf("failed to create material: %w", err)
	}

	return nil
}

// Get retrieves a material by ID
func (ref *MaterialRepository) Get(ctx context.Context, id string) (*domain.Material, error) {
	key := datastore.NameKey((&domain.Material{}).EntityKind(), id, nil)
	material := &domain.Material{}

	err := ref.client.client.Get(ctx, key, material)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, fmt.Errorf("material not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get material: %w", err)
	}

	material.ID = key.Name
	return material, nil
}

// List retrieves a list of materials
func (ref *MaterialRepository) List(ctx context.Context, limit, offset int) ([]*domain.Material, error) {
	query := datastore.NewQuery((&domain.Material{}).EntityKind()).
		Order("-createdAt").
		Limit(limit).
		Offset(offset)

	var materials []*domain.Material
	keys, err := ref.client.client.GetAll(ctx, query, &materials)
	if err != nil {
		return nil, fmt.Errorf("failed to list materials: %w", err)
	}

	for i, key := range keys {
		materials[i].ID = key.Name
	}

	return materials, nil
}

// ListByType retrieves materials filtered by type
func (ref *MaterialRepository) ListByType(ctx context.Context, materialType string, limit, offset int) ([]*domain.Material, error) {
	query := datastore.NewQuery((&domain.Material{}).EntityKind()).
		FilterField("type", "=", materialType).
		Order("-createdAt").
		Limit(limit).
		Offset(offset)

	var materials []*domain.Material
	keys, err := ref.client.client.GetAll(ctx, query, &materials)
	if err != nil {
		return nil, fmt.Errorf("failed to list materials by type: %w", err)
	}

	for i, key := range keys {
		materials[i].ID = key.Name
	}

	return materials, nil
}

// ListByLabel retrieves materials filtered by label
func (ref *MaterialRepository) ListByLabel(ctx context.Context, label string, limit, offset int) ([]*domain.Material, error) {
	query := datastore.NewQuery((&domain.Material{}).EntityKind()).
		FilterField("label", "=", label).
		Order("-createdAt").
		Limit(limit).
		Offset(offset)

	var materials []*domain.Material
	keys, err := ref.client.client.GetAll(ctx, query, &materials)
	if err != nil {
		return nil, fmt.Errorf("failed to list materials by label: %w", err)
	}

	for i, key := range keys {
		materials[i].ID = key.Name
	}

	return materials, nil
}

// ListByDateRange retrieves materials created within a date range
func (ref *MaterialRepository) ListByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Material, error) {
	query := datastore.NewQuery((&domain.Material{}).EntityKind()).
		FilterField("createdAt", ">=", start).
		FilterField("createdAt", "<=", end).
		Order("createdAt")

	var materials []*domain.Material
	keys, err := ref.client.client.GetAll(ctx, query, &materials)
	if err != nil {
		return nil, fmt.Errorf("failed to list materials by date range: %w", err)
	}

	for i, key := range keys {
		materials[i].ID = key.Name
	}

	return materials, nil
}

// Update updates an existing material
func (ref *MaterialRepository) Update(ctx context.Context, material *domain.Material) error {
	key := datastore.NameKey(material.EntityKind(), material.ID, nil)

	_, err := ref.client.client.Put(ctx, key, material)
	if err != nil {
		return fmt.Errorf("failed to update material: %w", err)
	}

	return nil
}

// Delete deletes a material by ID
func (ref *MaterialRepository) Delete(ctx context.Context, id string) error {
	key := datastore.NameKey((&domain.Material{}).EntityKind(), id, nil)

	err := ref.client.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete material: %w", err)
	}

	return nil
}
