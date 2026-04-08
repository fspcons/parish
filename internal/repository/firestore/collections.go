package firestore

// Firestore collection IDs (stable; not the same as domain.EntityKind strings).
const (
	colUsers        = "users"
	colRoles        = "roles"
	colEvents       = "events"
	colMaterials    = "materials"
	colParishGroups = "parish_groups"
	colSchedules    = "schedules"
)

const scheduleSingletonDocID = "parish-schedule"
