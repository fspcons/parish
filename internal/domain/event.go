package domain

// Event represents a parish event
type Event struct {
	BaseEntity
	Title       string `json:"title" datastore:"title"`
	Description string `json:"description" datastore:"description,noindex"`
	ImgURL      string `json:"imgUrl" datastore:"imgUrl,noindex"`
	Date        string `json:"date" datastore:"date"`
	Location    string `json:"location" datastore:"location"`
	Origin      string `json:"origin" datastore:"origin"`
}

// NewEvent creates a new Event entity. Returns an error if validation fails.
func NewEvent(title, description, imgURL, date, location, origin, createdBy string) (*Event, error) {
	e := &Event{
		BaseEntity:  NewBaseEntity(createdBy),
		Title:       title,
		Description: description,
		ImgURL:      imgURL,
		Date:        date,
		Location:    location,
		Origin:      origin,
	}
	if err := e.Validate(); err != nil {
		return nil, err
	}
	return e, nil
}

// Validate checks that the event satisfies its invariants.
func (ref *Event) Validate() error {
	if ref.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

// Update applies new field values, validates, and updates the timestamp.
func (ref *Event) Update(title, description, imgURL, date, location, origin, updatedBy string) error {
	ref.Title = title
	ref.Description = description
	ref.ImgURL = imgURL
	ref.Date = date
	ref.Location = location
	ref.Origin = origin
	if err := ref.Validate(); err != nil {
		return err
	}
	ref.UpdateTimestamp(updatedBy)
	return nil
}

// EntityKind returns the Datastore kind for this entity.
func (ref *Event) EntityKind() string {
	return "Event"
}
