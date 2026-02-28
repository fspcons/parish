package datastore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/parish/internal/domain"
)

// EventRepository implements repository.EventRepository
type EventRepository struct {
	client *Client
}

// NewEventRepository creates a new EventRepository
func NewEventRepository(client *Client) *EventRepository {
	return &EventRepository{
		client: client,
	}
}

// Create creates a new event
func (ref *EventRepository) Create(ctx context.Context, event *domain.Event) error {
	key := datastore.NameKey(event.EntityKind(), event.ID, nil)

	_, err := ref.client.client.Put(ctx, key, event)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

// Get retrieves an event by ID
func (ref *EventRepository) Get(ctx context.Context, id string) (*domain.Event, error) {
	key := datastore.NameKey((&domain.Event{}).EntityKind(), id, nil)
	event := &domain.Event{}

	err := ref.client.client.Get(ctx, key, event)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, fmt.Errorf("event not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	event.ID = key.Name
	return event, nil
}

// List retrieves a list of events
func (ref *EventRepository) List(ctx context.Context, limit, offset int) ([]*domain.Event, error) {
	query := datastore.NewQuery((&domain.Event{}).EntityKind()).
		Order("-createdAt").
		Limit(limit).
		Offset(offset)

	var events []*domain.Event
	keys, err := ref.client.client.GetAll(ctx, query, &events)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	for i, key := range keys {
		events[i].ID = key.Name
	}

	return events, nil
}

// ListByDateRange retrieves events created within a date range
func (ref *EventRepository) ListByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Event, error) {
	query := datastore.NewQuery((&domain.Event{}).EntityKind()).
		FilterField("createdAt", ">=", start).
		FilterField("createdAt", "<=", end).
		Order("createdAt")

	var events []*domain.Event
	keys, err := ref.client.client.GetAll(ctx, query, &events)
	if err != nil {
		return nil, fmt.Errorf("failed to list events by date range: %w", err)
	}

	for i, key := range keys {
		events[i].ID = key.Name
	}

	return events, nil
}

// Update updates an existing event
func (ref *EventRepository) Update(ctx context.Context, event *domain.Event) error {
	key := datastore.NameKey(event.EntityKind(), event.ID, nil)

	_, err := ref.client.client.Put(ctx, key, event)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

// Delete deletes an event by ID
func (ref *EventRepository) Delete(ctx context.Context, id string) error {
	key := datastore.NameKey((&domain.Event{}).EntityKind(), id, nil)

	err := ref.client.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}
