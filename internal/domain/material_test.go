package domain

import (
	"errors"
	"testing"
)

func TestNewMaterial(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		matType string
		wantErr error
	}{
		{"videos", "Lecture", MaterialTypeVideos, nil},
		{"documents", "Handout", MaterialTypeDocuments, nil},
		{"empty title", "", MaterialTypeVideos, ErrTitleRequired},
		{"invalid type", "Title", "podcasts", ErrInvalidMaterialType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMaterial(tt.title, tt.matType, "desc", "http://url", "label", "admin")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if m.Title != tt.title {
					t.Errorf("expected title %q, got %q", tt.title, m.Title)
				}
				if m.Type != tt.matType {
					t.Errorf("expected type %q, got %q", tt.matType, m.Type)
				}
			}
		})
	}
}

func TestMaterialValidate(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		matType string
		wantErr error
	}{
		{"valid videos", "T", MaterialTypeVideos, nil},
		{"valid documents", "T", MaterialTypeDocuments, nil},
		{"empty title", "", MaterialTypeVideos, ErrTitleRequired},
		{"invalid type", "T", "audio", ErrInvalidMaterialType},
		{"empty title takes precedence", "", "bad", ErrTitleRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Material{Title: tt.title, Type: tt.matType}
			if err := m.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMaterialUpdate(t *testing.T) {
	tests := []struct {
		name     string
		newTitle string
		newType  string
		wantErr  error
	}{
		{"success", "New", MaterialTypeDocuments, nil},
		{"invalid type", "Title", "invalid", ErrInvalidMaterialType},
		{"empty title", "", MaterialTypeVideos, ErrTitleRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMaterial("Old", MaterialTypeVideos, "old desc", "http://old", "old-label", "admin")
			err := m.Update(tt.newTitle, tt.newType, "new desc", "http://new", "new-label", "editor")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if m.Title != tt.newTitle {
					t.Errorf("expected title %q, got %q", tt.newTitle, m.Title)
				}
				if m.UpdatedBy != "editor" {
					t.Errorf("expected updatedBy 'editor', got %q", m.UpdatedBy)
				}
			}
		})
	}
}

func TestMaterialEntityKind(t *testing.T) {
	m := &Material{}
	if kind := m.EntityKind(); kind != "Material" {
		t.Errorf("expected 'Material', got %q", kind)
	}
}
