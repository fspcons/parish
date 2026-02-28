package usecase

import (
	"context"
	"log/slog"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"
)

// UpdateScheduleInput holds the fields required to update a schedule.
type UpdateScheduleInput struct {
	Monday    string
	Tuesday   string
	Wednesday string
	Thursday  string
	Friday    string
	Saturday  string
	Sunday    string
	UpdatedBy string
}

// Schedule defines schedule use case operations
type Schedule interface {
	Get(ctx context.Context) (*domain.Schedule, error)
	Update(ctx context.Context, in UpdateScheduleInput) (*domain.Schedule, error)
}

type schedule struct {
	repo repository.ScheduleRepository
}

// NewSchedule creates a new schedule use case
func NewSchedule(repo repository.ScheduleRepository) Schedule {
	return &schedule{
		repo: repo,
	}
}

// Get retrieves the schedule
func (ref *schedule) Get(ctx context.Context) (*domain.Schedule, error) {
	s, err := ref.repo.Get(ctx)
	if err != nil {
		slog.Error("failed to get schedule", "error", err)
		return nil, domain.ErrInternalServerError
	}
	return s, nil
}

// Update retrieves the current schedule, applies the new day values, and persists it.
func (ref *schedule) Update(ctx context.Context, in UpdateScheduleInput) (*domain.Schedule, error) {
	s, err := ref.repo.Get(ctx)
	if err != nil {
		slog.Error("failed to get schedule for update", "error", err)
		return nil, domain.ErrInternalServerError
	}

	s.UpdateDays(in.Monday, in.Tuesday, in.Wednesday, in.Thursday, in.Friday, in.Saturday, in.Sunday, in.UpdatedBy)

	if err := ref.repo.Put(ctx, s); err != nil {
		slog.Error("failed to persist schedule update", "error", err)
		return nil, domain.ErrInternalServerError
	}

	return s, nil
}
