package domain

import "golang.org/x/crypto/bcrypt"

// User represents a system user
type User struct {
	BaseEntity
	Email    string   `json:"email" datastore:"email"`
	Name     string   `json:"name" datastore:"name,noindex"`
	Password string   `json:"-" datastore:"password,noindex"` // Hashed password
	RoleIDs  []string `json:"roleIds" datastore:"roleIds"`
	Active   bool     `json:"active" datastore:"active"`
}

// NewUser creates a new User entity
func NewUser(email, name, hashedPassword string, roleIDs []string, createdBy string) *User {
	return &User{
		BaseEntity: NewBaseEntity(createdBy),
		Email:      email,
		Name:       name,
		Password:   hashedPassword,
		RoleIDs:    roleIDs,
		Active:     true,
	}
}

// Validate checks that the user satisfies its invariants.
func (ref *User) Validate() error {
	if ref.Email == "" {
		return ErrEmailRequired
	}
	if ref.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

// CheckPassword compares a plain-text password against the stored hash.
func (ref *User) CheckPassword(plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(ref.Password), []byte(plainPassword))
}

// IsActive returns whether the user account is active.
func (ref *User) IsActive() bool {
	return ref.Active
}

// Activate marks the user account as active.
func (ref *User) Activate(updatedBy string) {
	ref.Active = true
	ref.UpdateTimestamp(updatedBy)
}

// Deactivate marks the user account as inactive.
func (ref *User) Deactivate(updatedBy string) {
	ref.Active = false
	ref.UpdateTimestamp(updatedBy)
}

// EntityKind returns the Datastore kind for this entity.
func (ref *User) EntityKind() string {
	return "User"
}
