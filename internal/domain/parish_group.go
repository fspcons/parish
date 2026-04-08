package domain

// ParishGroup represents a parish group
type ParishGroup struct {
	BaseEntity
	Title       string `json:"title" firestore:"title"`
	Description string `json:"description" firestore:"description"`
	Manager     string `json:"manager" firestore:"manager"`
	Active      bool   `json:"active" firestore:"active"`
}

// NewParishGroup creates a new ParishGroup entity. Returns an error if validation fails.
func NewParishGroup(title, description, manager string, active bool, createdBy string) (*ParishGroup, error) {
	pg := &ParishGroup{
		BaseEntity:  NewBaseEntity(createdBy),
		Title:       title,
		Description: description,
		Manager:     manager,
		Active:      active,
	}
	if err := pg.Validate(); err != nil {
		return nil, err
	}
	return pg, nil
}

// Validate checks that the parish group satisfies its invariants.
func (ref *ParishGroup) Validate() error {
	if ref.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

// Update applies new field values, validates, and updates the timestamp.
func (ref *ParishGroup) Update(title, description, manager string, active bool, updatedBy string) error {
	ref.Title = title
	ref.Description = description
	ref.Manager = manager
	ref.Active = active
	if err := ref.Validate(); err != nil {
		return err
	}
	ref.UpdateTimestamp(updatedBy)
	return nil
}

// Activate marks the parish group as active.
func (ref *ParishGroup) Activate(updatedBy string) {
	ref.Active = true
	ref.UpdateTimestamp(updatedBy)
}

// Deactivate marks the parish group as inactive.
func (ref *ParishGroup) Deactivate(updatedBy string) {
	ref.Active = false
	ref.UpdateTimestamp(updatedBy)
}

// EntityKind returns the logical entity name (Firestore collection is "parish_groups").
func (ref ParishGroup) EntityKind() string {
	return "ParishGroup"
}

func (ref ParishGroup) SetID(id string) ParishGroup {
	ref.ID = id
	return ref
}
