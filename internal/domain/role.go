package domain

// Role represents a user role with permissions
type Role struct {
	BaseEntity
	Name        string       `json:"name" datastore:"name"`
	Description string       `json:"description" datastore:"description,noindex"`
	Permissions []Permission `json:"permissions" datastore:"permissions,noindex"`
}

// Permission represents a fine-grained permission
type Permission struct {
	Resource string `json:"resource"` // e.g., "schedule", "events", "parish_groups", "materials"
	Read     bool   `json:"read"`
	Write    bool   `json:"write"`
}

// NewRole creates a new Role entity
func NewRole(name, description string, permissions []Permission, createdBy string) *Role {
	return &Role{
		BaseEntity:  NewBaseEntity(createdBy),
		Name:        name,
		Description: description,
		Permissions: permissions,
	}
}

// HasPermission checks if the role has the specified permission for a resource.
func (ref *Role) HasPermission(resource string, write bool) bool {
	for _, perm := range ref.Permissions {
		if perm.Resource == resource {
			if write {
				return perm.Write
			}
			return perm.Read
		}
	}
	return false
}

// EntityKind returns the Datastore kind for this entity.
func (ref *Role) EntityKind() string {
	return "Role"
}
