package domain

// ScheduleResponse is the API representation of a parish schedule (day strings only).
type ScheduleResponse struct {
	Monday    string `json:"monday"`
	Tuesday   string `json:"tuesday"`
	Wednesday string `json:"wednesday"`
	Thursday  string `json:"thursday"`
	Friday    string `json:"friday"`
	Saturday  string `json:"saturday"`
	Sunday    string `json:"sunday"`
}

// ToResponse maps a Schedule entity to its API shape.
func (s *Schedule) ToResponse() ScheduleResponse {
	if s == nil {
		return ScheduleResponse{}
	}
	return ScheduleResponse{
		Monday:    s.Monday,
		Tuesday:   s.Tuesday,
		Wednesday: s.Wednesday,
		Thursday:  s.Thursday,
		Friday:    s.Friday,
		Saturday:  s.Saturday,
		Sunday:    s.Sunday,
	}
}

// EventResponse is the API representation of an event.
type EventResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImgURL      string `json:"imgUrl"`
	Date        string `json:"date"`
	Location    string `json:"location"`
	Origin      string `json:"origin"`
}

// ToResponse maps an Event entity to its API shape.
func (e *Event) ToResponse() EventResponse {
	if e == nil {
		return EventResponse{}
	}
	return EventResponse{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		ImgURL:      e.ImgURL,
		Date:        e.Date,
		Location:    e.Location,
		Origin:      e.Origin,
	}
}

// ToEventResponses maps a slice of events; empty input yields an empty JSON array (not null).
func ToEventResponses(events []*Event) []EventResponse {
	out := make([]EventResponse, 0, len(events))
	for _, e := range events {
		if e == nil {
			continue
		}
		out = append(out, e.ToResponse())
	}
	return out
}

// ParishGroupResponse is the API representation of a parish group.
type ParishGroupResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Manager     string `json:"manager"`
	Active      bool   `json:"active"`
}

// ToResponse maps a ParishGroup entity to its API shape.
func (g *ParishGroup) ToResponse() ParishGroupResponse {
	if g == nil {
		return ParishGroupResponse{}
	}
	return ParishGroupResponse{
		ID:          g.ID,
		Title:       g.Title,
		Description: g.Description,
		Manager:     g.Manager,
		Active:      g.Active,
	}
}

// ToParishGroupResponses maps a slice of parish groups; empty input yields an empty JSON array (not null).
func ToParishGroupResponses(groups []*ParishGroup) []ParishGroupResponse {
	out := make([]ParishGroupResponse, 0, len(groups))
	for _, g := range groups {
		if g == nil {
			continue
		}
		out = append(out, g.ToResponse())
	}
	return out
}

// MaterialResponse is the API representation of a material.
type MaterialResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Label       string `json:"label"`
}

// ToResponse maps a Material entity to its API shape.
func (m *Material) ToResponse() MaterialResponse {
	if m == nil {
		return MaterialResponse{}
	}
	return MaterialResponse{
		ID:          m.ID,
		Title:       m.Title,
		Type:        m.Type,
		Description: m.Description,
		URL:         m.URL,
		Label:       m.Label,
	}
}

// ToMaterialResponses maps a slice of materials; empty input yields an empty JSON array (not null).
func ToMaterialResponses(materials []*Material) []MaterialResponse {
	out := make([]MaterialResponse, 0, len(materials))
	for _, m := range materials {
		if m == nil {
			continue
		}
		out = append(out, m.ToResponse())
	}
	return out
}

// RoleResponse is the API representation of a role.
type RoleResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}

// ToResponse maps a Role entity to its API shape.
func (r *Role) ToResponse() RoleResponse {
	if r == nil {
		return RoleResponse{Permissions: []Permission{}}
	}
	perms := r.Permissions
	if perms == nil {
		perms = []Permission{}
	}
	return RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: perms,
	}
}

// ToRoleResponses maps a slice of roles; empty input yields an empty JSON array (not null).
func ToRoleResponses(roles []*Role) []RoleResponse {
	out := make([]RoleResponse, 0, len(roles))
	for _, r := range roles {
		if r == nil {
			continue
		}
		out = append(out, r.ToResponse())
	}
	return out
}

// UserResponse is the API representation of a user (no password or audit metadata).
type UserResponse struct {
	ID                 string   `json:"id"`
	Email              string   `json:"email"`
	Name               string   `json:"name"`
	RoleIDs            []string `json:"roleIds"`
	Active             bool     `json:"active"`
	MustChangePassword bool     `json:"mustChangePassword"`
}

// ToResponse maps a User entity to its API shape.
func (u *User) ToResponse() UserResponse {
	if u == nil {
		return UserResponse{RoleIDs: []string{}}
	}
	roleIDs := u.RoleIDs
	if roleIDs == nil {
		roleIDs = []string{}
	}
	return UserResponse{
		ID:                 u.ID,
		Email:              u.Email,
		Name:               u.Name,
		RoleIDs:            roleIDs,
		Active:             u.Active,
		MustChangePassword: u.MustChangePassword,
	}
}
