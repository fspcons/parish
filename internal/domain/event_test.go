package domain

import (
	"errors"
	"testing"
)

func TestNewEvent(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"success", "Mass", nil},
		{"empty title", "", ErrTitleRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEvent(tt.title, "desc", "http://img.png", "2026-01-01", "Church", "parish", "admin")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if e.Title != tt.title {
					t.Errorf("expected title %q, got %q", tt.title, e.Title)
				}
				if e.ID == "" {
					t.Error("expected ID to be set")
				}
				if e.CreatedBy != "admin" {
					t.Errorf("expected createdBy 'admin', got %q", e.CreatedBy)
				}
			}
		})
	}
}

func TestEventValidate(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"valid", "ok", nil},
		{"empty title", "", ErrTitleRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{Title: tt.title}
			if err := e.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventUpdate(t *testing.T) {
	tests := []struct {
		name     string
		newTitle string
		wantErr  error
	}{
		{"success", "New", nil},
		{"empty title", "", ErrTitleRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, _ := NewEvent("Old", "desc", "", "", "", "", "admin")
			err := e.Update(tt.newTitle, "new desc", "http://new.png", "2026-06-01", "Hall", "diocese", "editor")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if e.Title != tt.newTitle {
					t.Errorf("expected title %q, got %q", tt.newTitle, e.Title)
				}
				if e.UpdatedBy != "editor" {
					t.Errorf("expected updatedBy 'editor', got %q", e.UpdatedBy)
				}
			}
		})
	}
}

func TestEventEntityKind(t *testing.T) {
	e := &Event{}
	if kind := e.EntityKind(); kind != "Event" {
		t.Errorf("expected 'Event', got %q", kind)
	}
}
