package firestore

import (
	"context"
	"fmt"

	gcfs "cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/parish/internal/domain"
)

// RoleRepository implements repository.RoleRepository.
type RoleRepository struct {
	store *Store
}

// NewRoleRepository creates a RoleRepository.
func NewRoleRepository(store *Store) *RoleRepository {
	return &RoleRepository{store: store}
}

func (r *RoleRepository) col() *gcfs.CollectionRef {
	return r.store.Collection(colRoles)
}

// Create creates a role document.
func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	_, err := r.col().Doc(role.ID).Create(ctx, role)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return fmt.Errorf("role already exists: %s", role.ID)
		}
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// Get returns a role by ID.
func (r *RoleRepository) Get(ctx context.Context, id string) (*domain.Role, error) {
	snap, err := r.col().Doc(id).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("role not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	var role domain.Role
	if err := snap.DataTo(&role); err != nil {
		return nil, fmt.Errorf("failed to decode role: %w", err)
	}
	role.ID = id
	return &role, nil
}

// GetByName finds a role by name.
func (r *RoleRepository) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	q := r.col().Where("name", "==", name).Limit(1)
	iter := q.Documents(ctx)

	roles, err := scanDocuments[domain.Role](iter)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}
	if len(roles) == 0 {
		return nil, fmt.Errorf("role not found with name: %s", name)
	}

	return roles[0], nil
}

// List returns roles ordered by createdAt descending.
func (r *RoleRepository) List(ctx context.Context, limit, offset int) ([]*domain.Role, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	q := r.col().OrderBy("createdAt", gcfs.Desc).Limit(limit).Offset(offset)
	iter := q.Documents(ctx)

	return scanDocuments[domain.Role](iter)
}

// GetMultiple loads roles by ID (nil entries for missing IDs).
func (r *RoleRepository) GetMultiple(ctx context.Context, ids []string) ([]*domain.Role, error) {
	if len(ids) == 0 {
		return []*domain.Role{}, nil
	}

	refs := make([]*gcfs.DocumentRef, len(ids))
	for i, id := range ids {
		refs[i] = r.col().Doc(id)
	}

	snaps, err := r.store.GetAll(ctx, refs)
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple roles: %w", err)
	}

	roles := make([]*domain.Role, len(ids))
	for i, snap := range snaps {
		if snap == nil || !snap.Exists() {
			roles[i] = nil
			continue
		}
		var role domain.Role
		if err := snap.DataTo(&role); err != nil {
			return nil, fmt.Errorf("failed to decode role %s: %w", ids[i], err)
		}
		role.ID = ids[i]
		roles[i] = &role
	}
	return roles, nil
}

// Update saves a role.
func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	_, err := r.col().Doc(role.ID).Set(ctx, role)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

// Delete removes a role.
func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.col().Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}
