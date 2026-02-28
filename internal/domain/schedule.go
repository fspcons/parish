package domain

// Schedule represents the parish schedule (single row entity)
type Schedule struct {
	BaseEntity
	Monday    string `json:"monday" datastore:"monday,noindex"`
	Tuesday   string `json:"tuesday" datastore:"tuesday,noindex"`
	Wednesday string `json:"wednesday" datastore:"wednesday,noindex"`
	Thursday  string `json:"thursday" datastore:"thursday,noindex"`
	Friday    string `json:"friday" datastore:"friday,noindex"`
	Saturday  string `json:"saturday" datastore:"saturday,noindex"`
	Sunday    string `json:"sunday" datastore:"sunday,noindex"`
}

// NewSchedule creates a new Schedule entity
func NewSchedule(createdBy string) *Schedule {
	return &Schedule{
		BaseEntity: NewBaseEntity(createdBy),
	}
}

// UpdateDays applies new values to all day fields and updates the timestamp.
func (ref *Schedule) UpdateDays(monday, tuesday, wednesday, thursday, friday, saturday, sunday, updatedBy string) {
	ref.Monday = monday
	ref.Tuesday = tuesday
	ref.Wednesday = wednesday
	ref.Thursday = thursday
	ref.Friday = friday
	ref.Saturday = saturday
	ref.Sunday = sunday
	ref.UpdateTimestamp(updatedBy)
}

// EntityKind returns the Datastore kind for this entity.
func (ref *Schedule) EntityKind() string {
	return "Schedule"
}
