package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"
)

func TestMaterialCreate(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		matType   string
		createErr error
		wantErr   error
	}{
		{"success", "Lecture", domain.MaterialTypeVideos, nil, nil},
		{"empty title", "", domain.MaterialTypeVideos, nil, domain.ErrTitleRequired},
		{"invalid type", "Lecture", "podcasts", nil, domain.ErrInvalidMaterialType},
		{"repo error", "Lecture", domain.MaterialTypeVideos, errors.New("db"), domain.ErrInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.MaterialRepositoryMock{
				CreateFunc: func(_ context.Context, _ *domain.Material) error { return tt.createErr },
			}
			uc := NewMaterial(repo)
			mat, err := uc.Create(context.Background(), CreateMaterialInput{
				Title:       tt.title,
				Type:        tt.matType,
				Description: "desc",
				URL:         "http://url",
				Label:       "label",
				CreatedBy:   "admin",
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && mat.Title != tt.title {
				t.Errorf("expected title %q, got %q", tt.title, mat.Title)
			}
		})
	}
}

func TestMaterialGet(t *testing.T) {
	tests := []struct {
		name    string
		getFunc func(context.Context, string) (*domain.Material, error)
		wantErr error
	}{
		{
			"success",
			func(_ context.Context, _ string) (*domain.Material, error) {
				m := &domain.Material{}
				m.ID = "abc"
				return m, nil
			},
			nil,
		},
		{
			"not found",
			func(_ context.Context, _ string) (*domain.Material, error) { return nil, errors.New("not found") },
			domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.MaterialRepositoryMock{GetFunc: tt.getFunc}
			uc := NewMaterial(repo)
			_, err := uc.Get(context.Background(), "abc")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMaterialList(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		listFunc func(context.Context, int, int) ([]*domain.Material, error)
		wantErr  error
		wantLen  int
	}{
		{
			"success with default limit", 0,
			func(_ context.Context, limit, _ int) ([]*domain.Material, error) {
				if limit != 20 {
					return nil, errors.New("expected default limit 20")
				}
				return []*domain.Material{{}, {}}, nil
			},
			nil, 2,
		},
		{
			"repo error", 10,
			func(_ context.Context, _, _ int) ([]*domain.Material, error) { return nil, errors.New("db") },
			domain.ErrInternalServerError, 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.MaterialRepositoryMock{ListFunc: tt.listFunc}
			uc := NewMaterial(repo)
			mats, err := uc.List(context.Background(), tt.limit, 0)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && len(mats) != tt.wantLen {
				t.Errorf("expected %d, got %d", tt.wantLen, len(mats))
			}
		})
	}
}

func TestMaterialListByType(t *testing.T) {
	tests := []struct {
		name     string
		listFunc func(context.Context, string, int, int) ([]*domain.Material, error)
		wantErr  error
	}{
		{
			"success",
			func(_ context.Context, _ string, _, _ int) ([]*domain.Material, error) {
				return []*domain.Material{{}}, nil
			},
			nil,
		},
		{
			"repo error",
			func(_ context.Context, _ string, _, _ int) ([]*domain.Material, error) { return nil, errors.New("db") },
			domain.ErrInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.MaterialRepositoryMock{ListByTypeFunc: tt.listFunc}
			uc := NewMaterial(repo)
			_, err := uc.ListByType(context.Background(), domain.MaterialTypeVideos, 0, 0)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMaterialListByLabel(t *testing.T) {
	tests := []struct {
		name     string
		listFunc func(context.Context, string, int, int) ([]*domain.Material, error)
		wantErr  error
	}{
		{
			"success",
			func(_ context.Context, _ string, _, _ int) ([]*domain.Material, error) {
				return []*domain.Material{{}}, nil
			},
			nil,
		},
		{
			"repo error",
			func(_ context.Context, _ string, _, _ int) ([]*domain.Material, error) { return nil, errors.New("db") },
			domain.ErrInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.MaterialRepositoryMock{ListByLabelFunc: tt.listFunc}
			uc := NewMaterial(repo)
			_, err := uc.ListByLabel(context.Background(), "catechism", 0, 0)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMaterialUpdate(t *testing.T) {
	tests := []struct {
		name      string
		newTitle  string
		newType   string
		getFunc   func(context.Context, string) (*domain.Material, error)
		updateErr error
		wantErr   error
	}{
		{
			"success", "New", domain.MaterialTypeDocuments,
			func(_ context.Context, _ string) (*domain.Material, error) {
				m := &domain.Material{Type: domain.MaterialTypeVideos}
				m.ID = "abc"
				m.Title = "Old"
				return m, nil
			},
			nil, nil,
		},
		{
			"not found", "New", domain.MaterialTypeVideos,
			func(_ context.Context, _ string) (*domain.Material, error) { return nil, errors.New("not found") },
			nil, domain.ErrNotFound,
		},
		{
			"validation error", "", domain.MaterialTypeVideos,
			func(_ context.Context, _ string) (*domain.Material, error) {
				m := &domain.Material{Type: domain.MaterialTypeVideos}
				m.ID = "abc"
				m.Title = "Old"
				return m, nil
			},
			nil, domain.ErrTitleRequired,
		},
		{
			"repo update error", "New", domain.MaterialTypeVideos,
			func(_ context.Context, _ string) (*domain.Material, error) {
				m := &domain.Material{Type: domain.MaterialTypeVideos}
				m.ID = "abc"
				m.Title = "Old"
				return m, nil
			},
			errors.New("db"), domain.ErrInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository.MaterialRepositoryMock{
				GetFunc:    tt.getFunc,
				UpdateFunc: func(_ context.Context, _ *domain.Material) error { return tt.updateErr },
			}
			uc := NewMaterial(repo)
			_, err := uc.Update(context.Background(), UpdateMaterialInput{
				ID:          "abc",
				Title:       tt.newTitle,
				Type:        tt.newType,
				Description: "desc",
				URL:         "url",
				Label:       "label",
				UpdatedBy:   "editor",
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMaterialDelete(t *testing.T) {
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
			repo := &repository.MaterialRepositoryMock{
				DeleteFunc: func(_ context.Context, _ string) error { return tt.deleteErr },
			}
			uc := NewMaterial(repo)
			err := uc.Delete(context.Background(), "abc")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
