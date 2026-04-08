package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Entity[T any] interface {
	EntityKind() string
	SetID(id string) T
}

// BaseEntity contains common fields for all entities
type BaseEntity struct {
	ID        string    `json:"id" firestore:"-"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`
	CreatedBy string    `json:"createdBy" firestore:"createdBy"`
	UpdatedBy string    `json:"updatedBy" firestore:"updatedBy"`
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

// IsEmpty checks if a string is empty
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
