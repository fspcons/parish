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

// EventRepository implements repository.EventRepository.
type EventRepository struct {
	store *Store
}

// NewEventRepository creates an EventRepository.
func NewEventRepository(store *Store) *EventRepository {
	return &EventRepository{store: store}
}

func (r *EventRepository) col() *gcfs.CollectionRef {
	return r.store.Collection(colEvents)
}

// Create creates an event.
func (r *EventRepository) Create(ctx context.Context, event *domain.Event) error {
	_, err := r.col().Doc(event.ID).Create(ctx, event)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return fmt.Errorf("event already exists: %s", event.ID)
		}
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

// Get returns an event by ID.
func (r *EventRepository) Get(ctx context.Context, id string) (*domain.Event, error) {
	snap, err := r.col().Doc(id).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("event not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}
	var e domain.Event
	if err := snap.DataTo(&e); err != nil {
		return nil, fmt.Errorf("failed to decode event: %w", err)
	}
	e.ID = id
	return &e, nil
}

// List lists events by createdAt descending.
func (r *EventRepository) List(ctx context.Context, limit, offset int) ([]*domain.Event, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	q := r.col().OrderBy("createdAt", gcfs.Desc).Limit(limit).Offset(offset)
	iter := q.Documents(ctx)

	return scanDocuments[domain.Event](iter)
}

// ListByDateRange lists events whose createdAt is in [start, end].
func (r *EventRepository) ListByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Event, error) {
	q := r.col().
		Where("createdAt", ">=", start).
		Where("createdAt", "<=", end).
		OrderBy("createdAt", gcfs.Asc)
	iter := q.Documents(ctx)

	return scanDocuments[domain.Event](iter)
}

// Update saves an event.
func (r *EventRepository) Update(ctx context.Context, event *domain.Event) error {
	_, err := r.col().Doc(event.ID).Set(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	return nil
}

// Delete removes an event.
func (r *EventRepository) Delete(ctx context.Context, id string) error {
	_, err := r.col().Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}
