package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"
)

func TestParishGroupCreate(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		createErr error
		wantErr   error
	}{
		{"success", "Youth", nil, nil},
		{"empty title", "", nil, domain.ErrTitleRequired},
		{"repo error", "Youth", errors.New("db"), domain.ErrInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.ParishGroupRepositoryMock{
				CreateFunc: func(_ context.Context, _ *domain.ParishGroup) error { return tt.createErr },
			}
			uc := NewParishGroup(repo)
			g, err := uc.Create(context.Background(), CreateParishGroupInput{
				Title:       tt.title,
				Description: "desc",
				Manager:     "mgr",
				Active:      true,
				CreatedBy:   "admin",
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && g.Title != tt.title {
				t.Errorf("expected title %q, got %q", tt.title, g.Title)
			}
		})
	}
}

func TestParishGroupGet(t *testing.T) {
	tests := []struct {
		name    string
		getFunc func(context.Context, string) (*domain.ParishGroup, error)
		wantErr error
	}{
		{
			"success",
			func(_ context.Context, _ string) (*domain.ParishGroup, error) {
				g := &domain.ParishGroup{}
				g.ID = "abc"
				return g, nil
			},
			nil,
		},
		{
			"not found",
			func(_ context.Context, _ string) (*domain.ParishGroup, error) { return nil, errors.New("not found") },
			domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.ParishGroupRepositoryMock{GetFunc: tt.getFunc}
			uc := NewParishGroup(repo)
			_, err := uc.Get(context.Background(), "abc")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestParishGroupList(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		listFunc func(context.Context, int, int) ([]*domain.ParishGroup, error)
		wantErr  error
		wantLen  int
	}{
		{
			"success with default limit", 0,
			func(_ context.Context, limit, _ int) ([]*domain.ParishGroup, error) {
				if limit != 20 {
					return nil, errors.New("expected default limit 20")
				}
				return []*domain.ParishGroup{{}, {}}, nil
			},
			nil, 2,
		},
		{
			"repo error", 10,
			func(_ context.Context, _, _ int) ([]*domain.ParishGroup, error) { return nil, errors.New("db") },
			domain.ErrInternalServerError, 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.ParishGroupRepositoryMock{ListFunc: tt.listFunc}
			uc := NewParishGroup(repo)
			groups, err := uc.List(context.Background(), tt.limit, 0)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && len(groups) != tt.wantLen {
				t.Errorf("expected %d, got %d", tt.wantLen, len(groups))
			}
		})
	}
}

func TestParishGroupUpdate(t *testing.T) {
	tests := []struct {
		name      string
		newTitle  string
		getFunc   func(context.Context, string) (*domain.ParishGroup, error)
		updateErr error
		wantErr   error
	}{
		{
			"success", "New",
			func(_ context.Context, _ string) (*domain.ParishGroup, error) {
				g := &domain.ParishGroup{}
				g.ID = "abc"
				g.Title = "Old"
				return g, nil
			},
			nil, nil,
		},
		{
			"not found", "New",
			func(_ context.Context, _ string) (*domain.ParishGroup, error) { return nil, errors.New("not found") },
			nil, domain.ErrNotFound,
		},
		{
			"validation error", "",
			func(_ context.Context, _ string) (*domain.ParishGroup, error) {
				g := &domain.ParishGroup{}
				g.ID = "abc"
				g.Title = "Old"
				return g, nil
			},
			nil, domain.ErrTitleRequired,
		},
		{
			"repo update error", "New",
			func(_ context.Context, _ string) (*domain.ParishGroup, error) {
				g := &domain.ParishGroup{}
				g.ID = "abc"
				g.Title = "Old"
				return g, nil
			},
			errors.New("db"), domain.ErrInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.ParishGroupRepositoryMock{
				GetFunc:    tt.getFunc,
				UpdateFunc: func(_ context.Context, _ *domain.ParishGroup) error { return tt.updateErr },
			}
			uc := NewParishGroup(repo)
			_, err := uc.Update(context.Background(), UpdateParishGroupInput{
				ID:          "abc",
				Title:       tt.newTitle,
				Description: "desc",
				Manager:     "mgr",
				Active:      true,
				UpdatedBy:   "editor",
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestParishGroupDelete(t *testing.T) {
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
			repo := &repository.ParishGroupRepositoryMock{
				DeleteFunc: func(_ context.Context, _ string) error { return tt.deleteErr },
			}
			uc := NewParishGroup(repo)
			err := uc.Delete(context.Background(), "abc")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
