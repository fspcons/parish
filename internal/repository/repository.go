package repository

//go:generate moq -out mocks.go . ScheduleRepository ParishGroupRepository EventRepository MaterialRepository UserRepository RoleRepository

import (
	"context"
	"time"

	"github.com/parish/internal/domain"
)

type (
	// ScheduleRepository defines operations for schedule persistence
	ScheduleRepository interface {
		Get(ctx context.Context) (*domain.Schedule, error)
		Put(ctx context.Context, schedule *domain.Schedule) error
	}

	// ParishGroupRepository defines operations for parish group persistence
	ParishGroupRepository interface {
		Create(ctx context.Context, group *domain.ParishGroup) error
		Get(ctx context.Context, id string) (*domain.ParishGroup, error)
		List(ctx context.Context, limit, offset int) ([]*domain.ParishGroup, error)
		ListByDateRange(ctx context.Context, start, end time.Time) ([]*domain.ParishGroup, error)
		Update(ctx context.Context, group *domain.ParishGroup) error
		Delete(ctx context.Context, id string) error
	}

	// EventRepository defines operations for event persistence
	EventRepository interface {
		Create(ctx context.Context, event *domain.Event) error
		Get(ctx context.Context, id string) (*domain.Event, error)
		List(ctx context.Context, limit, offset int) ([]*domain.Event, error)
		ListByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Event, error)
		Update(ctx context.Context, event *domain.Event) error
		Delete(ctx context.Context, id string) error
	}

	// MaterialRepository defines operations for material persistence
	MaterialRepository interface {
		Create(ctx context.Context, material *domain.Material) error
		Get(ctx context.Context, id string) (*domain.Material, error)
		List(ctx context.Context, limit, offset int) ([]*domain.Material, error)
		ListByType(ctx context.Context, materialType string, limit, offset int) ([]*domain.Material, error)
		ListByLabel(ctx context.Context, label string, limit, offset int) ([]*domain.Material, error)
		ListByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Material, error)
		Update(ctx context.Context, material *domain.Material) error
		Delete(ctx context.Context, id string) error
	}

	// UserRepository defines operations for user persistence
	UserRepository interface {
		Create(ctx context.Context, user *domain.User) error
		Get(ctx context.Context, id string) (*domain.User, error)
		GetByEmail(ctx context.Context, email string) (*domain.User, error)
		List(ctx context.Context, limit, offset int) ([]*domain.User, error)
		Update(ctx context.Context, user *domain.User) error
		Delete(ctx context.Context, id string) error
	}

	// RoleRepository defines operations for role persistence
	RoleRepository interface {
		Create(ctx context.Context, role *domain.Role) error
		Get(ctx context.Context, id string) (*domain.Role, error)
		GetByName(ctx context.Context, name string) (*domain.Role, error)
		List(ctx context.Context, limit, offset int) ([]*domain.Role, error)
		GetMultiple(ctx context.Context, ids []string) ([]*domain.Role, error)
		Update(ctx context.Context, role *domain.Role) error
		Delete(ctx context.Context, id string) error
	}
)
