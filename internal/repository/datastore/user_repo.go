package datastore

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/parish/internal/domain"
)

// UserRepository implements repository.UserRepository
type UserRepository struct {
	client *Client
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(client *Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// Create creates a new user
func (ref *UserRepository) Create(ctx context.Context, user *domain.User) error {
	key := datastore.NameKey(user.EntityKind(), user.ID, nil)

	_, err := ref.client.client.Put(ctx, key, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Get retrieves a user by ID
func (ref *UserRepository) Get(ctx context.Context, id string) (*domain.User, error) {
	key := datastore.NameKey((&domain.User{}).EntityKind(), id, nil)
	user := &domain.User{}

	err := ref.client.client.Get(ctx, key, user)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.ID = key.Name
	return user, nil
}

// GetByEmail retrieves a user by email
func (ref *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := datastore.NewQuery((&domain.User{}).EntityKind()).
		FilterField("email", "=", email).
		Limit(1)

	var users []*domain.User
	keys, err := ref.client.client.GetAll(ctx, query, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found with email: %s", email)
	}

	users[0].ID = keys[0].Name
	return users[0], nil
}

// List retrieves a list of users
func (ref *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := datastore.NewQuery((&domain.User{}).EntityKind()).
		Order("-createdAt").
		Limit(limit).
		Offset(offset)

	var users []*domain.User
	keys, err := ref.client.client.GetAll(ctx, query, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	for i, key := range keys {
		users[i].ID = key.Name
	}

	return users, nil
}

// Update updates an existing user
func (ref *UserRepository) Update(ctx context.Context, user *domain.User) error {
	key := datastore.NameKey(user.EntityKind(), user.ID, nil)

	_, err := ref.client.client.Put(ctx, key, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user by ID
func (ref *UserRepository) Delete(ctx context.Context, id string) error {
	key := datastore.NameKey((&domain.User{}).EntityKind(), id, nil)

	err := ref.client.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
