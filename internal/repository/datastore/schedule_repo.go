package datastore

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/parish/internal/domain"
)

// ScheduleRepository implements repository.ScheduleRepository
type ScheduleRepository struct {
	client *Client
}

// NewScheduleRepository creates a new ScheduleRepository
func NewScheduleRepository(client *Client) *ScheduleRepository {
	return &ScheduleRepository{
		client: client,
	}
}

const scheduleKey = "parish-schedule"

// Get retrieves the schedule (single entity)
func (ref *ScheduleRepository) Get(ctx context.Context) (*domain.Schedule, error) {
	key := datastore.NameKey((&domain.Schedule{}).EntityKind(), scheduleKey, nil)
	schedule := &domain.Schedule{}

	err := ref.client.client.Get(ctx, key, schedule)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			// Return empty schedule if it doesn't exist
			return &domain.Schedule{
				BaseEntity: domain.NewBaseEntity("system"),
			}, nil
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	schedule.ID = key.Name
	return schedule, nil
}

// Put creates or updates the schedule
func (ref *ScheduleRepository) Put(ctx context.Context, schedule *domain.Schedule) error {
	key := datastore.NameKey(schedule.EntityKind(), scheduleKey, nil)

	_, err := ref.client.client.Put(ctx, key, schedule)
	if err != nil {
		return fmt.Errorf("failed to save schedule: %w", err)
	}

	schedule.ID = key.Name
	return nil
}
