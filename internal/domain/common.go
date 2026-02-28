package domain

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity contains common fields for all entities
type BaseEntity struct {
	ID        string    `json:"id" datastore:"-"`
	CreatedAt time.Time `json:"createdAt" datastore:"createdAt,noindex"`
	UpdatedAt time.Time `json:"updatedAt" datastore:"updatedAt,noindex"`
	CreatedBy string    `json:"createdBy" datastore:"createdBy,noindex"`
	UpdatedBy string    `json:"updatedBy" datastore:"updatedBy,noindex"`
}

// NewBaseEntity creates a new base entity with UUID and timestamps
func NewBaseEntity(createdBy string) BaseEntity {
	now := time.Now().UTC()
	return BaseEntity{
		ID:        uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}
}

// UpdateTimestamp updates the UpdatedAt and UpdatedBy fields
func (ref *BaseEntity) UpdateTimestamp(updatedBy string) {
	ref.UpdatedAt = time.Now().UTC()
	ref.UpdatedBy = updatedBy
}
