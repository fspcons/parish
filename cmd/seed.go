package main

import (
	"context"
	"log/slog"

	"github.com/parish/internal/domain"
	"github.com/parish/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// seedAdmin idempotently creates an admin role and admin user on startup.
// It is a no-op when ADMIN_EMAIL or ADMIN_PASSWORD are empty, or when the
// admin user already exists.
func seedAdmin(ctx context.Context, cfg *Config, userRepo repository.UserRepository, roleRepo repository.RoleRepository) {
	if cfg.AdminEmail == "" || cfg.AdminPassword == "" {
		return
	}

	// ----- idempotent admin role -----
	const adminRoleName = "admin"
	adminRole, err := roleRepo.GetByName(ctx, adminRoleName)
	if err != nil {
		// Role doesn't exist yet — create it.
		adminRole = domain.NewRole(adminRoleName, "System administrator with full permissions", []domain.Permission{
			{Resource: "schedule", Read: true, Write: true},
			{Resource: "events", Read: true, Write: true},
			{Resource: "parish_groups", Read: true, Write: true},
			{Resource: "materials", Read: true, Write: true},
			{Resource: "roles", Read: true, Write: true},
		}, "system")

		if err := roleRepo.Create(ctx, adminRole); err != nil {
			slog.Error("seed: failed to create admin role", "error", err)
			return
		}
		slog.Info("seed: admin role created", "roleID", adminRole.ID)
	} else {
		slog.Info("seed: admin role already exists", "roleID", adminRole.ID)
	}

	// ----- idempotent admin user -----
	existingUser, err := userRepo.GetByEmail(ctx, cfg.AdminEmail)
	if err == nil && existingUser != nil {
		slog.Info("seed: admin user already exists", "email", cfg.AdminEmail)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("seed: failed to hash admin password", "error", err)
		return
	}

	adminUser := domain.NewUser(cfg.AdminEmail, "Admin", string(hashedPassword), []string{adminRole.ID}, "system")
	if err := userRepo.Create(ctx, adminUser); err != nil {
		slog.Error("seed: failed to create admin user", "error", err)
		return
	}

	slog.Info("seed: admin user created", "email", cfg.AdminEmail, "userID", adminUser.ID)
}
