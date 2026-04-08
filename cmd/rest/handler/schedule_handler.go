package handler

import (
	"encoding/json"
	"net/http"

	"github.com/parish/cmd/rest/middleware"
	"github.com/parish/internal/usecase"
)

// ScheduleHandler handles schedule requests
type ScheduleHandler struct {
	scheduleUseCase usecase.Schedule
}

// NewScheduleHandler creates a new schedule handler
func NewScheduleHandler(scheduleUseCase usecase.Schedule) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleUseCase: scheduleUseCase,
	}
}

// ScheduleRequest represents a schedule update request
type ScheduleRequest struct {
	Monday    string `json:"monday"`
	Tuesday   string `json:"tuesday"`
	Wednesday string `json:"wednesday"`
	Thursday  string `json:"thursday"`
	Friday    string `json:"friday"`
	Saturday  string `json:"saturday"`
	Sunday    string `json:"sunday"`
}

// Get retrieves the schedule
func (ref *ScheduleHandler) Get(w http.ResponseWriter, r *http.Request) {
	s, err := ref.scheduleUseCase.Get(r.Context())
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, s.ToResponse(), "")
}

// Update updates the schedule
func (ref *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	s, err := ref.scheduleUseCase.Update(r.Context(), usecase.UpdateScheduleInput{
		Monday:    req.Monday,
		Tuesday:   req.Tuesday,
		Wednesday: req.Wednesday,
		Thursday:  req.Thursday,
		Friday:    req.Friday,
		Saturday:  req.Saturday,
		Sunday:    req.Sunday,
		UpdatedBy: user.ID,
	})
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, s.ToResponse(), "Schedule updated successfully")
}
