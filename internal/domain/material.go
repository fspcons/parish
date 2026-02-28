package domain

const (
	MaterialTypeVideos    = "videos"
	MaterialTypeDocuments = "documents"
)

// Material represents a parish material (video or document)
type Material struct {
	BaseEntity
	Title       string `json:"title" datastore:"title"`
	Type        string `json:"type" datastore:"type"` // "videos" or "documents"
	Description string `json:"description" datastore:"description,noindex"`
	URL         string `json:"url" datastore:"url,noindex"`
	Label       string `json:"label" datastore:"label"` // Hierarchical label separated by colons
}

// NewMaterial creates a new Material entity. Returns an error if validation fails.
func NewMaterial(title, materialType, description, url, label, createdBy string) (*Material, error) {
	m := &Material{
		BaseEntity:  NewBaseEntity(createdBy),
		Title:       title,
		Type:        materialType,
		Description: description,
		URL:         url,
		Label:       label,
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

// Validate checks that the material satisfies its invariants.
func (ref *Material) Validate() error {
	if ref.Title == "" {
		return ErrTitleRequired
	}
	if ref.Type != MaterialTypeVideos && ref.Type != MaterialTypeDocuments {
		return ErrInvalidMaterialType
	}
	return nil
}

// Update applies new field values, validates, and updates the timestamp.
func (ref *Material) Update(title, materialType, description, url, label, updatedBy string) error {
	ref.Title = title
	ref.Type = materialType
	ref.Description = description
	ref.URL = url
	ref.Label = label
	if err := ref.Validate(); err != nil {
		return err
	}
	ref.UpdateTimestamp(updatedBy)
	return nil
}

// EntityKind returns the Datastore kind for this entity.
func (ref *Material) EntityKind() string {
	return "Material"
}
