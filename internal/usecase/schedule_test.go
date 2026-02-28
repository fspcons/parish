package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"
)

func TestScheduleGet(t *testing.T) {
	tests := []struct {
		name    string
		getFunc func(context.Context) (*domain.Schedule, error)
		wantErr error
	}{
		{
			"success",
			func(_ context.Context) (*domain.Schedule, error) { return domain.NewSchedule("system"), nil },
			nil,
		},
		{
			"repo error",
			func(_ context.Context) (*domain.Schedule, error) { return nil, errors.New("db") },
			domain.ErrInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.ScheduleRepositoryMock{GetFunc: tt.getFunc}
			uc := NewSchedule(repo)
			s, err := uc.Get(context.Background())
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && s == nil {
				t.Error("expected schedule, got nil")
			}
		})
	}
}

func TestScheduleUpdate(t *testing.T) {
	tests := []struct {
		name         string
		getFunc      func(context.Context) (*domain.Schedule, error)
		createOrUpFn func(context.Context, *domain.Schedule) error
		wantErr      error
	}{
		{
			"success",
			func(_ context.Context) (*domain.Schedule, error) { return domain.NewSchedule("system"), nil },
			func(_ context.Context, _ *domain.Schedule) error { return nil },
			nil,
		},
		{
			"get error",
			func(_ context.Context) (*domain.Schedule, error) { return nil, errors.New("db") },
			func(_ context.Context, _ *domain.Schedule) error { return nil },
			domain.ErrInternalServerError,
		},
		{
			"persist error",
			func(_ context.Context) (*domain.Schedule, error) { return domain.NewSchedule("system"), nil },
			func(_ context.Context, _ *domain.Schedule) error { return errors.New("db") },
			domain.ErrInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.ScheduleRepositoryMock{
				GetFunc: tt.getFunc,
				PutFunc: tt.createOrUpFn,
			}
			uc := NewSchedule(repo)
			s, err := uc.Update(context.Background(), UpdateScheduleInput{
				Monday:    "Mon",
				Tuesday:   "Tue",
				Wednesday: "Wed",
				Thursday:  "Thu",
				Friday:    "Fri",
				Saturday:  "Sat",
				Sunday:    "Sun",
				UpdatedBy: "editor",
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if s.Monday != "Mon" {
					t.Errorf("expected Monday 'Mon', got %q", s.Monday)
				}
				if s.UpdatedBy != "editor" {
					t.Errorf("expected updatedBy 'editor', got %q", s.UpdatedBy)
				}
			}
		})
	}
}
