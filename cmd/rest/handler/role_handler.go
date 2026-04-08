package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/parish/cmd/rest/middleware"
	"github.com/parish/internal/domain"
	"github.com/parish/internal/usecase"
)

// RoleHandler handles role requests
type RoleHandler struct {
	roleUseCase usecase.Role
	authUseCase usecase.Auth
}

// NewRoleHandler creates a new role handler
func NewRoleHandler(roleUseCase usecase.Role, authUseCase usecase.Auth) *RoleHandler {
	return &RoleHandler{
		roleUseCase: roleUseCase,
		authUseCase: authUseCase,
	}
}

// Create creates a new role
func (ref *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req RoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role, err := ref.roleUseCase.Create(r.Context(), usecase.CreateRoleInput{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		CreatedBy:   user.ID,
	})
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusCreated, role.ToResponse(), "Role created successfully")
}

// Get retrieves a role by ID
func (ref *RoleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	role, err := ref.roleUseCase.Get(r.Context(), id)
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, role.ToResponse(), "")
}

// List retrieves a list of roles
func (ref *RoleHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	roles, err := ref.roleUseCase.List(r.Context(), limit, offset)
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, domain.ToRoleResponses(roles), "")
}

// Update updates a role
func (ref *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req RoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role, err := ref.roleUseCase.Update(r.Context(), usecase.UpdateRoleInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		UpdatedBy:   user.ID,
	})
	if err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, role.ToResponse(), "Role updated successfully")
}

// Delete deletes a role
func (ref *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := ref.roleUseCase.Delete(r.Context(), id); err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, nil, "Role deleted successfully")
}

// AssignRoles assigns role IDs to a user
func (ref *RoleHandler) AssignRoles(w http.ResponseWriter, r *http.Request) {
	caller := middleware.GetUserFromContext(r.Context())
	if caller == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req AssignRolesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := ref.authUseCase.AssignRoles(r.Context(), userID, req.RoleIDs, caller.ID); err != nil {
		HandleDomainError(w, err)
		return
	}

	RespondSuccess(w, http.StatusOK, nil, "Roles assigned successfully")
}
