package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/parish/cmd/rest/middleware"
	"github.com/parish/internal/domain"
	"github.com/parish/internal/usecase"
)

// MaterialHandler handles material requests
type MaterialHandler struct {
	useCase usecase.Material
}

// NewMaterialHandler creates a new material handler
func NewMaterialHandler(useCase usecase.Material) *MaterialHandler {
	return &MaterialHandler{
		useCase: useCase,
	}
}

// Create creates a new material
func (ref *MaterialHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req MaterialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mat, err := ref.useCase.Create(r.Context(), usecase.CreateMaterialInput{
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
		URL:         req.URL,
		Label:       req.Label,
		CreatedBy:   user.ID,
	})
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusCreated, mat.ToResponse(), "Material created successfully")
}

// Get retrieves a material by ID
func (ref *MaterialHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mat, err := ref.useCase.Get(r.Context(), id)
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, mat.ToResponse(), "")
}

// List retrieves a list of materials
func (ref *MaterialHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	materialType := r.URL.Query().Get("type")
	label := r.URL.Query().Get("label")

	var materials []*domain.Material
	var err error

	if materialType != "" {
		materials, err = ref.useCase.ListByType(r.Context(), materialType, limit, offset)
	} else if label != "" {
		materials, err = ref.useCase.ListByLabel(r.Context(), label, limit, offset)
	} else {
		materials, err = ref.useCase.List(r.Context(), limit, offset)
	}

	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, domain.ToMaterialResponses(materials), "")
}

// Update updates a material
func (ref *MaterialHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req MaterialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mat, err := ref.useCase.Update(r.Context(), usecase.UpdateMaterialInput{
		ID:          id,
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
		URL:         req.URL,
		Label:       req.Label,
		UpdatedBy:   user.ID,
	})
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, mat.ToResponse(), "Material updated successfully")
}

// Delete deletes a material
func (ref *MaterialHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	RespondSuccess(w, http.StatusOK, nil, "Material deleted successfully")
}
