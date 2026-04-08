package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/parish/cmd/rest/middleware"
	"github.com/parish/internal/domain"
	"github.com/parish/internal/usecase"
)

// ParishGroupHandler handles parish group requests
type ParishGroupHandler struct {
	useCase usecase.ParishGroup
}

// NewParishGroupHandler creates a new parish group handler
func NewParishGroupHandler(useCase usecase.ParishGroup) *ParishGroupHandler {
	return &ParishGroupHandler{
		useCase: useCase,
	}
}

// Create creates a new parish group
func (ref *ParishGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ParishGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	group, err := ref.useCase.Create(r.Context(), usecase.CreateParishGroupInput{
		Title:       req.Title,
		Description: req.Description,
		Manager:     req.Manager,
		Active:      req.Active,
		CreatedBy:   user.ID,
	})
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusCreated, group.ToResponse(), "Parish group created successfully")
}

// Get retrieves a parish group by ID
func (ref *ParishGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	group, err := ref.useCase.Get(r.Context(), id)
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, group.ToResponse(), "")
}

// List retrieves a list of parish groups
func (ref *ParishGroupHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	groups, err := ref.useCase.List(r.Context(), limit, offset)
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, domain.ToParishGroupResponses(groups), "")
}

// Update updates a parish group
func (ref *ParishGroupHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req ParishGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	group, err := ref.useCase.Update(r.Context(), usecase.UpdateParishGroupInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Manager:     req.Manager,
		Active:      req.Active,
		UpdatedBy:   user.ID,
	})
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, group.ToResponse(), "Parish group updated successfully")
}

// Delete deletes a parish group
func (ref *ParishGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err := ref.useCase.Delete(r.Context(), id)
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, nil, "Parish group deleted successfully")
}
