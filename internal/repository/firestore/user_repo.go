package firestore

import (
	"context"
	"fmt"

	gcfs "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/parish/internal/domain"
)

// UserRepository implements repository.UserRepository.
type UserRepository struct {
	store *Store
}

// NewUserRepository creates a UserRepository.
func NewUserRepository(store *Store) *UserRepository {
	return &UserRepository{store: store}
}

func (r *UserRepository) col() *gcfs.CollectionRef {
	return r.store.Collection(colUsers)
}

// Create creates a new user document.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.col().Doc(user.ID).Create(ctx, user)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return fmt.Errorf("user already exists: %s", user.ID)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Get retrieves a user by ID.
func (r *UserRepository) Get(ctx context.Context, id string) (*domain.User, error) {
	snap, err := r.col().Doc(id).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	var u domain.User
	if err := snap.DataTo(&u); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}
	u.ID = id
	return &u, nil
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := r.col().Where("email", "==", email).Limit(1)
	iter := q.Documents(ctx)
	

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("user not found with email: %s", email)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	var u domain.User
	if err := doc.DataTo(&u); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}
	u.ID = doc.Ref.ID
	return &u, nil
}

// List lists users ordered by createdAt descending.
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	q := r.col().OrderBy("createdAt", gcfs.Desc).Limit(limit).Offset(offset)
	iter := q.Documents(ctx)
	

	return scanDocuments[domain.User](iter)
}

// Update replaces a user document.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := r.col().Doc(user.ID).Set(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete removes a user.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.col().Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
