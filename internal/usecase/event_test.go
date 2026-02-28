package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"
)

func TestEventCreate(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		createErr error
		wantErr   error
	}{
		{"success", "Mass", nil, nil},
		{"validation error", "", nil, domain.ErrTitleRequired},
		{"repo error", "Mass", errors.New("db"), domain.ErrInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.EventRepositoryMock{
				CreateFunc: func(_ context.Context, _ *domain.Event) error { return tt.createErr },
			}
			uc := NewEvent(repo)
			evt, err := uc.Create(context.Background(), CreateEventInput{
				Title:       tt.title,
				Description: "desc",
				Date:        "2026-01-01",
				Location:    "Church",
				Origin:      "parish",
				CreatedBy:   "admin",
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if evt.Title != tt.title {
					t.Errorf("expected title %q, got %q", tt.title, evt.Title)
				}
				if len(repo.CreateCalls()) != 1 {
					t.Errorf("expected 1 Create call, got %d", len(repo.CreateCalls()))
				}
			}
		})
	}
}

func TestEventGet(t *testing.T) {
	existing := &domain.Event{}
	existing.ID = "abc"
	existing.Title = "Mass"

	tests := []struct {
		name    string
		id      string
		getFunc func(context.Context, string) (*domain.Event, error)
		wantErr error
	}{
		{
			"success", "abc",
			func(_ context.Context, id string) (*domain.Event, error) { return existing, nil },
			nil,
		},
		{
			"not found", "missing",
			func(_ context.Context, id string) (*domain.Event, error) { return nil, errors.New("not found") },
			domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.EventRepositoryMock{
				GetFunc: tt.getFunc,
			}
			uc := NewEvent(repo)
			_, err := uc.Get(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventList(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		listFunc func(context.Context, int, int) ([]*domain.Event, error)
		wantErr  error
		wantLen  int
	}{
		{
			"success with default limit", 0,
			func(_ context.Context, limit, offset int) ([]*domain.Event, error) {
				if limit != 20 {
					return nil, errors.New("expected default limit 20")
				}
				return []*domain.Event{{}, {}}, nil
			},
			nil, 2,
		},
		{
			"repo error", 10,
			func(_ context.Context, _, _ int) ([]*domain.Event, error) { return nil, errors.New("db") },
			domain.ErrInternalServerError, 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.EventRepositoryMock{
				ListFunc: tt.listFunc,
			}
			uc := NewEvent(repo)
			events, err := uc.List(context.Background(), tt.limit, 0)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && len(events) != tt.wantLen {
				t.Errorf("expected %d events, got %d", tt.wantLen, len(events))
			}
		})
	}
}

func TestEventUpdate(t *testing.T) {
	existing := &domain.Event{}
	existing.ID = "abc"
	existing.Title = "Old"

	tests := []struct {
		name      string
		id        string
		newTitle  string
		getFunc   func(context.Context, string) (*domain.Event, error)
		updateErr error
		wantErr   error
	}{
		{
			"success", "abc", "New",
			func(_ context.Context, _ string) (*domain.Event, error) {
				e := &domain.Event{}
				e.ID = "abc"
				e.Title = "Old"
				return e, nil
			},
			nil, nil,
		},
		{
			"not found", "missing", "New",
			func(_ context.Context, _ string) (*domain.Event, error) { return nil, errors.New("not found") },
			nil, domain.ErrNotFound,
		},
		{
			"validation error", "abc", "",
			func(_ context.Context, _ string) (*domain.Event, error) {
				e := &domain.Event{}
				e.ID = "abc"
				e.Title = "Old"
				return e, nil
			},
			nil, domain.ErrTitleRequired,
		},
		{
			"repo update error", "abc", "New",
			func(_ context.Context, _ string) (*domain.Event, error) {
				e := &domain.Event{}
				e.ID = "abc"
				e.Title = "Old"
				return e, nil
			},
			errors.New("db"), domain.ErrInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.EventRepositoryMock{
				GetFunc:    tt.getFunc,
				UpdateFunc: func(_ context.Context, _ *domain.Event) error { return tt.updateErr },
			}
			uc := NewEvent(repo)
			evt, err := uc.Update(context.Background(), UpdateEventInput{
				ID:          tt.id,
				Title:       tt.newTitle,
				Description: "desc",
				UpdatedBy:   "editor",
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && evt.Title != tt.newTitle {
				t.Errorf("expected title %q, got %q", tt.newTitle, evt.Title)
			}
		})
	}
}

func TestEventDelete(t *testing.T) {
	tests := []struct {
		name      string
		deleteErr error
		wantErr   error
	}{
		{"success", nil, nil},
		{"repo error", errors.New("db"), domain.ErrInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.EventRepositoryMock{
				DeleteFunc: func(_ context.Context, _ string) error { return tt.deleteErr },
			}
			uc := NewEvent(repo)
			err := uc.Delete(context.Background(), "abc")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
