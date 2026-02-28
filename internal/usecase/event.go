package usecase

import (
	"context"
	"log/slog"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"
)

// CreateEventInput holds the fields required to create a new event.
type CreateEventInput struct {
	Title       string
	Description string
	ImgURL      string
	Date        string
	Location    string
	Origin      string
	CreatedBy   string
}

// UpdateEventInput holds the fields required to update an existing event.
type UpdateEventInput struct {
	ID          string
	Title       string
	Description string
	ImgURL      string
	Date        string
	Location    string
	Origin      string
	UpdatedBy   string
}

// Event defines event use case operations
type Event interface {
	Create(ctx context.Context, in CreateEventInput) (*domain.Event, error)
	Get(ctx context.Context, id string) (*domain.Event, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Event, error)
	Update(ctx context.Context, in UpdateEventInput) (*domain.Event, error)
	Delete(ctx context.Context, id string) error
}

type event struct {
	repo repository.EventRepository
}

// NewEvent creates a new event use case
func NewEvent(repo repository.EventRepository) Event {
	return &event{
		repo: repo,
	}
}

// Create creates a new event
func (ref *event) Create(ctx context.Context, in CreateEventInput) (*domain.Event, error) {
	evt, err := domain.NewEvent(in.Title, in.Description, in.ImgURL, in.Date, in.Location, in.Origin, in.CreatedBy)
	if err != nil {
		slog.Error("failed to create event entity", "error", err, "title", in.Title)
		return nil, err
	}

	if err := ref.repo.Create(ctx, evt); err != nil {
		slog.Error("failed to persist event", "error", err, "eventID", evt.ID)
		return nil, domain.ErrInternalServerError
	}

	return evt, nil
}

// Get retrieves an event by ID
func (ref *event) Get(ctx context.Context, id string) (*domain.Event, error) {
	evt, err := ref.repo.Get(ctx, id)
	if err != nil {
		slog.Error("failed to get event", "error", err, "eventID", id)
		return nil, domain.ErrNotFound
	}
	return evt, nil
}

// List retrieves a list of events
func (ref *event) List(ctx context.Context, limit, offset int) ([]*domain.Event, error) {
	if limit <= 0 {
		limit = 20
	}
	events, err := ref.repo.List(ctx, limit, offset)
	if err != nil {
		slog.Error("failed to list events", "error", err, "limit", limit, "offset", offset)
		return nil, domain.ErrInternalServerError
	}
	return events, nil
}

// Update updates an existing event
func (ref *event) Update(ctx context.Context, in UpdateEventInput) (*domain.Event, error) {
	evt, err := ref.repo.Get(ctx, in.ID)
	if err != nil {
		slog.Error("failed to get event for update", "error", err, "eventID", in.ID)
		return nil, domain.ErrNotFound
	}

	if err := evt.Update(in.Title, in.Description, in.ImgURL, in.Date, in.Location, in.Origin, in.UpdatedBy); err != nil {
		slog.Error("failed to update event entity", "error", err, "eventID", in.ID)
		return nil, err
	}

	if err := ref.repo.Update(ctx, evt); err != nil {
		slog.Error("failed to persist event update", "error", err, "eventID", in.ID)
		return nil, domain.ErrInternalServerError
	}

	return evt, nil
}

// Delete deletes an event
func (ref *event) Delete(ctx context.Context, id string) error {
	if err := ref.repo.Delete(ctx, id); err != nil {
		slog.Error("failed to delete event", "error", err, "eventID", id)
		return domain.ErrInternalServerError
	}
	return nil
}
