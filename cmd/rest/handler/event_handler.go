package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/parish/cmd/rest/middleware"
	"github.com/parish/internal/domain"
	"github.com/parish/internal/usecase"
)

// EventHandler handles event requests
type EventHandler struct {
	useCase usecase.Event
}

// NewEventHandler creates a new event handler
func NewEventHandler(useCase usecase.Event) *EventHandler {
	return &EventHandler{
		useCase: useCase,
	}
}

// Create creates a new event
func (ref *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	evt, err := ref.useCase.Create(r.Context(), usecase.CreateEventInput{
		Title:       req.Title,
		Description: req.Description,
		ImgURL:      req.ImgURL,
		Date:        req.Date,
		Location:    req.Location,
		Origin:      req.Origin,
		CreatedBy:   user.ID,
	})
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusCreated, evt.ToResponse(), "Event created successfully")
}

// Get retrieves an event by ID
func (ref *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	evt, err := ref.useCase.Get(r.Context(), id)
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, evt.ToResponse(), "")
}

// List retrieves a list of events
func (ref *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	events, err := ref.useCase.List(r.Context(), limit, offset)
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, domain.ToEventResponses(events), "")
}

// Update updates an event
func (ref *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	evt, err := ref.useCase.Update(r.Context(), usecase.UpdateEventInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		ImgURL:      req.ImgURL,
		Date:        req.Date,
		Location:    req.Location,
		Origin:      req.Origin,
		UpdatedBy:   user.ID,
	})
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, evt.ToResponse(), "Event updated successfully")
}

// Delete deletes an event
func (ref *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := ref.useCase.Delete(r.Context(), id); err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, nil, "Event deleted successfully")
}
