package domain

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func testHashPassword(t *testing.T, plain string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(h)
}

func TestNewUser(t *testing.T) {
	u := NewUser("test@example.com", "Test", "hashed", []string{"role1"}, "admin")
	if u.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", u.Email)
	}
	if !u.Active {
		t.Error("expected new user to be active")
	}
	if u.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestUserValidate(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  error
	}{
		{"valid", "a@b.com", "hashed", nil},
		{"empty email", "", "hashed", ErrEmailRequired},
		{"empty password", "a@b.com", "", ErrPasswordRequired},
		{"both empty returns email error first", "", "", ErrEmailRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Email: tt.email, Password: tt.password}
			if err := u.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserCheckPassword(t *testing.T) {
	hashed := testHashPassword(t, "secret123")
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"correct", "secret123", false},
		{"wrong", "wrongpassword", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Password: hashed}
			err := u.CheckPassword(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestUserIsActive(t *testing.T) {
	tests := []struct {
		name   string
		active bool
		want   bool
	}{
		{"active", true, true},
		{"inactive", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Active: tt.active}
			if got := u.IsActive(); got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestUserActivateDeactivate(t *testing.T) {
	tests := []struct {
		name       string
		action     string
		wantActive bool
	}{
		{"activate", "activate", true},
		{"deactivate", "deactivate", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUser("a@b.com", "A", "hashed", nil, "admin")
			if tt.action == "activate" {
				u.Active = false
				u.Activate("editor")
			} else {
				u.Deactivate("editor")
			}
			if u.Active != tt.wantActive {
				t.Errorf("expected active %v, got %v", tt.wantActive, u.Active)
			}
			if u.UpdatedBy != "editor" {
				t.Errorf("expected updatedBy 'editor', got %q", u.UpdatedBy)
			}
		})
	}
}

func TestUserEntityKind(t *testing.T) {
	u := &User{}
	if kind := u.EntityKind(); kind != "User" {
		t.Errorf("expected 'User', got %q", kind)
	}
}
