package datastore

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/parish/internal/domain"
)

// RoleRepository implements repository.RoleRepository
type RoleRepository struct {
	client *Client
}

// NewRoleRepository creates a new RoleRepository
func NewRoleRepository(client *Client) *RoleRepository {
	return &RoleRepository{
		client: client,
	}
}

// Create creates a new role
func (ref *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	key := datastore.NameKey(role.EntityKind(), role.ID, nil)

	_, err := ref.client.client.Put(ctx, key, role)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// Get retrieves a role by ID
func (ref *RoleRepository) Get(ctx context.Context, id string) (*domain.Role, error) {
	key := datastore.NameKey((&domain.Role{}).EntityKind(), id, nil)
	role := &domain.Role{}

	err := ref.client.client.Get(ctx, key, role)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, fmt.Errorf("role not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	role.ID = key.Name
	return role, nil
}

// GetByName retrieves a role by name
func (ref *RoleRepository) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	query := datastore.NewQuery((&domain.Role{}).EntityKind()).
		FilterField("name", "=", name).
		Limit(1)

	var roles []*domain.Role
	keys, err := ref.client.client.GetAll(ctx, query, &roles)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}

	if len(roles) == 0 {
		return nil, fmt.Errorf("role not found with name: %s", name)
	}

	roles[0].ID = keys[0].Name
	return roles[0], nil
}

// List retrieves a list of roles
func (ref *RoleRepository) List(ctx context.Context, limit, offset int) ([]*domain.Role, error) {
	query := datastore.NewQuery((&domain.Role{}).EntityKind()).
		Order("-createdAt").
		Limit(limit).
		Offset(offset)

	var roles []*domain.Role
	keys, err := ref.client.client.GetAll(ctx, query, &roles)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	for i, key := range keys {
		roles[i].ID = key.Name
	}

	return roles, nil
}

// GetMultiple retrieves multiple roles by their IDs
func (ref *RoleRepository) GetMultiple(ctx context.Context, ids []string) ([]*domain.Role, error) {
	if len(ids) == 0 {
		return []*domain.Role{}, nil
	}

	keys := make([]*datastore.Key, len(ids))
	for i, id := range ids {
		keys[i] = datastore.NameKey((&domain.Role{}).EntityKind(), id, nil)
	}

	roles := make([]*domain.Role, len(ids))
	err := ref.client.client.GetMulti(ctx, keys, roles)
	if err != nil {
		// Handle partial errors
		if multiErr, ok := err.(datastore.MultiError); ok {
			for i, e := range multiErr {
				if e != nil && e != datastore.ErrNoSuchEntity {
					return nil, fmt.Errorf("failed to get role %s: %w", ids[i], e)
				}
			}
		} else {
			return nil, fmt.Errorf("failed to get multiple roles: %w", err)
		}
	}

	for i, key := range keys {
		if roles[i] != nil {
			roles[i].ID = key.Name
		}
	}

	return roles, nil
}

// Update updates an existing role
func (ref *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	key := datastore.NameKey(role.EntityKind(), role.ID, nil)

	_, err := ref.client.client.Put(ctx, key, role)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

// Delete deletes a role by ID
func (ref *RoleRepository) Delete(ctx context.Context, id string) error {
	key := datastore.NameKey((&domain.Role{}).EntityKind(), id, nil)

	err := ref.client.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}
