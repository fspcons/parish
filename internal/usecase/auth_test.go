package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/parish/internal/cache"
	"github.com/parish/internal/domain"
	"github.com/parish/internal/email"
	"github.com/parish/internal/repository"
)

// noopMailer satisfies email.Sender for tests that do not assert on e-mail.
type noopMailer struct{}

func (noopMailer) SendPasswordReset(context.Context, string, string) error { return nil }

var _ email.Sender = noopMailer{}

func newCacheMock() *cache.CacheMock {
	store := make(map[string]string)
	return &cache.CacheMock{
		SetFunc: func(_ context.Context, key string, value string, _ time.Duration) error {
			store[key] = value
			return nil
		},
		GetFunc: func(_ context.Context, key string) (string, error) {
			v, ok := store[key]
			if !ok {
				return "", cache.ErrNotFound
			}
			return v, nil
		},
		DelFunc: func(_ context.Context, key string) error {
			delete(store, key)
			return nil
		},
	}
}

// newAuthMocks creates typical mock repos for auth tests.
func newAuthMocks() (*repository.UserRepositoryMock, *repository.RoleRepositoryMock) {
	users := make(map[string]*domain.User)
	byEmail := make(map[string]*domain.User)

	userRepo := &repository.UserRepositoryMock{
		CreateFunc: func(_ context.Context, u *domain.User) error {
			users[u.ID] = u
			byEmail[u.Email] = u
			return nil
		},
		GetFunc: func(_ context.Context, id string) (*domain.User, error) {
			u, ok := users[id]
			if !ok {
				return nil, errors.New("not found")
			}
			return u, nil
		},
		GetByEmailFunc: func(_ context.Context, email string) (*domain.User, error) {
			u, ok := byEmail[email]
			if !ok {
				return nil, errors.New("not found")
			}
			return u, nil
		},
		UpdateFunc: func(_ context.Context, u *domain.User) error {
			users[u.ID] = u
			byEmail[u.Email] = u
			return nil
		},
	}

	roles := make(map[string]*domain.Role)
	roleRepo := &repository.RoleRepositoryMock{
		CreateFunc: func(_ context.Context, r *domain.Role) error {
			roles[r.ID] = r
			return nil
		},
		GetMultipleFunc: func(_ context.Context, ids []string) ([]*domain.Role, error) {
			result := make([]*domain.Role, 0, len(ids))
			for _, id := range ids {
				if r, ok := roles[id]; ok {
					result = append(result, r)
				}
			}
			return result, nil
		},
	}

	return userRepo, roleRepo
}

func TestAuthRegister(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  error
	}{
		{"success", "test@example.com", "password123", nil},
		{"empty email", "", "password123", domain.ErrEmailRequired},
		{"empty password", "test@example.com", "", domain.ErrPasswordRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, roleRepo := newAuthMocks()
			uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})
			user, err := uc.Register(context.Background(), RegisterInput{
				Email:     tt.email,
				Name:      "Test",
				Password:  tt.password,
				CreatedBy: "system",
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if user.Email != tt.email {
					t.Errorf("expected email %q, got %q", tt.email, user.Email)
				}
				if !user.Active {
					t.Error("expected user to be active")
				}
			}
		})
	}
}

func TestAuthRegister_DuplicateEmail(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})

	_, _ = uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", CreatedBy: "system"})
	_, err := uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test2", Password: "password456", CreatedBy: "system"})
	if !errors.Is(err, domain.ErrUserAlreadyExists) {
		t.Errorf("got %v, want %v", err, domain.ErrUserAlreadyExists)
	}
}

func TestAuthRegister_RepoError(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	userRepo.CreateFunc = func(_ context.Context, _ *domain.User) error { return errors.New("db") }
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})

	_, err := uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", CreatedBy: "system"})
	if !errors.Is(err, domain.ErrInternalServerError) {
		t.Errorf("got %v, want %v", err, domain.ErrInternalServerError)
	}
}

func TestAuthLogin(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})

	_, _ = uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", CreatedBy: "system"})

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  error
	}{
		{"success", "test@example.com", "password123", nil},
		{"wrong password", "test@example.com", "wrong", domain.ErrInvalidCredentials},
		{"user not found", "unknown@example.com", "password", domain.ErrInvalidCredentials},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, token, err := uc.Login(context.Background(), tt.email, tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if user.Email != tt.email {
					t.Errorf("expected email %q, got %q", tt.email, user.Email)
				}
				if token == "" {
					t.Error("expected non-empty token")
				}
			}
		})
	}
}

func TestAuthLogin_InactiveUser(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})

	user, _ := uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", CreatedBy: "system"})
	user.Deactivate("admin")
	_ = userRepo.Update(context.Background(), user)

	_, _, err := uc.Login(context.Background(), "test@example.com", "password123")
	if !errors.Is(err, domain.ErrUserInactive) {
		t.Errorf("got %v, want %v", err, domain.ErrUserInactive)
	}
}

func TestAuthValidateToken(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})

	_, _ = uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", CreatedBy: "system"})
	_, token, _ := uc.Login(context.Background(), "test@example.com", "password123")

	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{"valid token", token, nil},
		{"invalid token", "invalid-token", domain.ErrInvalidCredentials},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := uc.ValidateToken(context.Background(), tt.token)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && user.Email != "test@example.com" {
				t.Errorf("expected email 'test@example.com', got %q", user.Email)
			}
		})
	}
}

func TestAuthCheckPermission(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})

	// Create roles
	role := domain.NewRole("admin", "Admin", []domain.Permission{
		{Resource: "events", Read: true, Write: true},
	}, "system")
	_ = roleRepo.Create(context.Background(), role)

	user, _ := uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", RoleIDs: []string{role.ID}, CreatedBy: "system"})

	tests := []struct {
		name     string
		resource string
		write    bool
		want     bool
	}{
		{"write granted", "events", true, true},
		{"read granted", "events", false, true},
		{"no matching resource", "schedule", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := uc.CheckPermission(context.Background(), user.ID, tt.resource, tt.write)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if allowed != tt.want {
				t.Errorf("got %v, want %v", allowed, tt.want)
			}
		})
	}
}

func TestAuthGetUserPermissions_MergesRoles(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})

	role1 := domain.NewRole("reader", "Reader", []domain.Permission{
		{Resource: "events", Read: true, Write: false},
	}, "system")
	role2 := domain.NewRole("writer", "Writer", []domain.Permission{
		{Resource: "events", Read: false, Write: true},
	}, "system")
	_ = roleRepo.Create(context.Background(), role1)
	_ = roleRepo.Create(context.Background(), role2)

	user, _ := uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", RoleIDs: []string{role1.ID, role2.ID}, CreatedBy: "system"})

	permissions, err := uc.GetUserPermissions(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(permissions) != 1 {
		t.Fatalf("expected 1 merged permission, got %d", len(permissions))
	}
	if !permissions[0].Read || !permissions[0].Write {
		t.Errorf("expected merged read+write, got read=%v write=%v", permissions[0].Read, permissions[0].Write)
	}
}

type mailRecorder struct {
	calls []struct {
		To    string
		Plain string
	}
}

func (m *mailRecorder) SendPasswordReset(_ context.Context, toEmail, plain string) error {
	m.calls = append(m.calls, struct {
		To    string
		Plain string
	}{To: toEmail, Plain: plain})
	return nil
}

var _ email.Sender = (*mailRecorder)(nil)

func TestRequestPasswordReset_EmptyEmail(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	m := &mailRecorder{}
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), m)
	uc.RequestPasswordReset(context.Background(), "   ")
	if len(m.calls) != 0 {
		t.Errorf("expected no email, got %d", len(m.calls))
	}
}

func TestRequestPasswordReset_UserNotFound(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	m := &mailRecorder{}
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), m)
	uc.RequestPasswordReset(context.Background(), "nobody@example.com")
	if len(m.calls) != 0 {
		t.Errorf("expected no email, got %d", len(m.calls))
	}
}

func TestRequestPasswordReset_InactiveUser(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})
	user, _ := uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", CreatedBy: "system"})
	user.Deactivate("admin")
	_ = userRepo.Update(context.Background(), user)

	m := &mailRecorder{}
	uc2 := NewAuth(userRepo, roleRepo, newCacheMock(), m)
	uc2.RequestPasswordReset(context.Background(), "test@example.com")
	if len(m.calls) != 0 {
		t.Errorf("expected no email for inactive user, got %d", len(m.calls))
	}
}

func TestRequestPasswordReset_Success(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	m := &mailRecorder{}
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), m)
	_, _ = uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", CreatedBy: "system"})

	uc.RequestPasswordReset(context.Background(), "test@example.com")
	if len(m.calls) != 1 {
		t.Fatalf("expected 1 email, got %d", len(m.calls))
	}
	if m.calls[0].To != "test@example.com" {
		t.Errorf("expected email to test@example.com, got %q", m.calls[0].To)
	}
	if m.calls[0].Plain == "" {
		t.Error("expected non-empty temporary password")
	}
	if m.calls[0].Plain == "password123" {
		t.Error("temporary password should not equal old password")
	}

	u, _, err := uc.Login(context.Background(), "test@example.com", m.calls[0].Plain)
	if err != nil {
		t.Fatalf("login with temp password: %v", err)
	}
	if !u.MustChangePassword {
		t.Error("expected MustChangePassword true after reset")
	}
}

func TestChangePassword_RequiresNewPassword(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})
	user, _ := uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", CreatedBy: "system"})

	err := uc.ChangePassword(context.Background(), user.ID, "password123", "")
	if !errors.Is(err, domain.ErrPasswordRequired) {
		t.Errorf("got %v, want ErrPasswordRequired", err)
	}
}

func TestChangePassword_WrongCurrentPassword(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), noopMailer{})
	user, _ := uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", CreatedBy: "system"})

	err := uc.ChangePassword(context.Background(), user.ID, "wrong", "newpass")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("got %v, want ErrInvalidCredentials", err)
	}
}

func TestChangePassword_Success_ClearsMustChangePassword(t *testing.T) {
	userRepo, roleRepo := newAuthMocks()
	m := &mailRecorder{}
	uc := NewAuth(userRepo, roleRepo, newCacheMock(), m)
	u, _ := uc.Register(context.Background(), RegisterInput{Email: "test@example.com", Name: "Test", Password: "password123", CreatedBy: "system"})
	uc.RequestPasswordReset(context.Background(), "test@example.com")
	temp := m.calls[0].Plain

	err := uc.ChangePassword(context.Background(), u.ID, temp, "brandNewPassword123")
	if err != nil {
		t.Fatalf("ChangePassword: %v", err)
	}
	userAfter, _, err := uc.Login(context.Background(), "test@example.com", "brandNewPassword123")
	if err != nil {
		t.Fatalf("login with new password: %v", err)
	}
	if userAfter.MustChangePassword {
		t.Error("expected MustChangePassword false after change")
	}
}
