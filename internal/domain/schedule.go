package domain

// Schedule represents the parish schedule (single row entity)
type Schedule struct {
	BaseEntity
	Monday    string `json:"monday" firestore:"monday"`
	Tuesday   string `json:"tuesday" firestore:"tuesday"`
	Wednesday string `json:"wednesday" firestore:"wednesday"`
	Thursday  string `json:"thursday" firestore:"thursday"`
	Friday    string `json:"friday" firestore:"friday"`
	Saturday  string `json:"saturday" firestore:"saturday"`
	Sunday    string `json:"sunday" firestore:"sunday"`
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

// EntityKind returns the logical entity name (Firestore collection is "schedules").
func (ref Schedule) EntityKind() string {
	return "Schedule"
}

func (ref Schedule) SetID(id string) Schedule {
	ref.ID = id
	return ref
}
