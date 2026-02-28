package domain

import (
	"errors"
	"testing"
)

func TestNewParishGroup(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"success", "Youth", nil},
		{"empty title", "", ErrTitleRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg, err := NewParishGroup(tt.title, "desc", "mgr", true, "admin")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if pg.Title != tt.title {
					t.Errorf("expected title %q, got %q", tt.title, pg.Title)
				}
				if !pg.Active {
					t.Error("expected active to be true")
				}
			}
		})
	}
}

func TestParishGroupValidate(t *testing.T) {
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
			pg := &ParishGroup{Title: tt.title}
			if err := pg.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestParishGroupUpdate(t *testing.T) {
	tests := []struct {
		name     string
		newTitle string
		active   bool
		wantErr  error
	}{
		{"success", "New", false, nil},
		{"empty title", "", true, ErrTitleRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg, _ := NewParishGroup("Old", "old desc", "OldMgr", true, "admin")
			err := pg.Update(tt.newTitle, "new desc", "NewMgr", tt.active, "editor")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if pg.Title != tt.newTitle {
					t.Errorf("expected title %q, got %q", tt.newTitle, pg.Title)
				}
				if pg.Active != tt.active {
					t.Errorf("expected active %v, got %v", tt.active, pg.Active)
				}
				if pg.UpdatedBy != "editor" {
					t.Errorf("expected updatedBy 'editor', got %q", pg.UpdatedBy)
				}
			}
		})
	}
}

func TestParishGroupActivateDeactivate(t *testing.T) {
	tests := []struct {
		name       string
		startActive bool
		action     string
		wantActive bool
	}{
		{"activate inactive", false, "activate", true},
		{"deactivate active", true, "deactivate", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg, _ := NewParishGroup("Group", "desc", "Mgr", tt.startActive, "admin")
			if tt.action == "activate" {
				pg.Activate("editor")
			} else {
				pg.Deactivate("editor")
			}
			if pg.Active != tt.wantActive {
				t.Errorf("expected active %v, got %v", tt.wantActive, pg.Active)
			}
			if pg.UpdatedBy != "editor" {
				t.Errorf("expected updatedBy 'editor', got %q", pg.UpdatedBy)
			}
		})
	}
}

func TestParishGroupEntityKind(t *testing.T) {
	pg := &ParishGroup{}
	if kind := pg.EntityKind(); kind != "ParishGroup" {
		t.Errorf("expected 'ParishGroup', got %q", kind)
	}
}
