package domain

import "golang.org/x/crypto/bcrypt"

// User represents a system user
type User struct {
	BaseEntity
	Email              string   `json:"email" firestore:"email"`
	Name               string   `json:"name" firestore:"name"`
	Password           string   `json:"-" firestore:"password"` // Hashed password
	RoleIDs            []string `json:"roleIds" firestore:"roleIds"`
	Active             bool     `json:"active" firestore:"active"`
	MustChangePassword bool     `json:"mustChangePassword" firestore:"mustChangePassword"`
}

// NewUser creates a new User entity
func NewUser(email, name, hashedPassword string, roleIDs []string, createdBy string) *User {
	return &User{
		BaseEntity:         NewBaseEntity(createdBy),
		Email:              email,
		Name:               name,
		Password:           hashedPassword,
		RoleIDs:            roleIDs,
		Active:             true,
		MustChangePassword: false,
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

// ApplyTemporaryPassword sets a new bcrypt-hashed password and marks the account so the client must change password after login.
func (ref *User) ApplyTemporaryPassword(hashedPassword string, updatedBy string) {
	ref.Password = hashedPassword
	ref.MustChangePassword = true
	ref.UpdateTimestamp(updatedBy)
}

// SetPasswordFromUserChange replaces the password after the user submits a new one and clears the forced-change flag.
func (ref *User) SetPasswordFromUserChange(hashedPassword string, updatedBy string) {
	ref.Password = hashedPassword
	ref.MustChangePassword = false
	ref.UpdateTimestamp(updatedBy)
}

// EntityKind returns the logical entity name (Firestore collection is "users").
func (ref User) EntityKind() string {
	return "User"
}

func (ref User) SetID(id string) User {
	ref.ID = id
	return ref
}
