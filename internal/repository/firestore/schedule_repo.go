package firestore

import (
	"context"
	"fmt"

	gcfs "cloud.google.com/go/firestore"

	"github.com/parish/internal/domain"
)

// ScheduleRepository implements repository.ScheduleRepository.
type ScheduleRepository struct {
	store *Store
}

// NewScheduleRepository creates a ScheduleRepository.
func NewScheduleRepository(store *Store) *ScheduleRepository {
	return &ScheduleRepository{store: store}
}

func (r *ScheduleRepository) doc() *gcfs.DocumentRef {
	return r.store.Collection(colSchedules).Doc(scheduleSingletonDocID)
}

// Get returns the singleton schedule document or a default empty schedule.
func (r *ScheduleRepository) Get(ctx context.Context) (*domain.Schedule, error) {
	snap, err := r.doc().Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return &domain.Schedule{
				BaseEntity: domain.NewBaseEntity("system"),
			}, nil
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}
	var s domain.Schedule
	if err := snap.DataTo(&s); err != nil {
		return nil, fmt.Errorf("failed to decode schedule: %w", err)
	}
	s.ID = scheduleSingletonDocID
	return &s, nil
}

// Put saves the singleton schedule.
func (r *ScheduleRepository) Put(ctx context.Context, schedule *domain.Schedule) error {
	_, err := r.doc().Set(ctx, schedule)
	if err != nil {
		return fmt.Errorf("failed to save schedule: %w", err)
	}
	schedule.ID = scheduleSingletonDocID
	return nil
}
