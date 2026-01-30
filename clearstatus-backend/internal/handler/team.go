package handler

import (
	"net/http"

	"github.com/clearstatus/backend/internal/middleware"
	"github.com/clearstatus/backend/internal/model"
	"github.com/clearstatus/backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TeamHandler struct {
	org  *store.OrganizationStore
	team *store.TeamStore
}

func NewTeamHandler(org *store.OrganizationStore, team *store.TeamStore) *TeamHandler {
	return &TeamHandler{org: org, team: team}
}

func (h *TeamHandler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{teamID}", h.Get)
	r.Put("/{teamID}", h.Update)
	r.Delete("/{teamID}", h.Delete)
}

func (h *TeamHandler) ensureOrgMember(w http.ResponseWriter, r *http.Request, orgID string) bool {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	ok, _ := h.org.UserInOrg(r.Context(), orgID, user.ID)
	if !ok {
		Error(w, http.StatusNotFound, "organization not found")
		return false
	}
	return true
}

func (h *TeamHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	list, err := h.team.ListByOrgID(r.Context(), orgID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, list)
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &body); err != nil {
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Name == "" {
		Error(w, http.StatusBadRequest, "name required")
		return
	}
	t := &model.Team{ID: uuid.New().String(), OrgID: orgID, Name: body.Name}
	if err := h.team.Create(r.Context(), t); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, t)
}

func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	teamID := chi.URLParam(r, "teamID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	t, err := h.team.GetByID(r.Context(), teamID)
	if err != nil {
		Error(w, http.StatusNotFound, "team not found")
		return
	}
	if t.OrgID != orgID {
		Error(w, http.StatusNotFound, "team not found")
		return
	}
	JSON(w, http.StatusOK, t)
}

func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	teamID := chi.URLParam(r, "teamID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	t, err := h.team.GetByID(r.Context(), teamID)
	if err != nil || t.OrgID != orgID {
		Error(w, http.StatusNotFound, "team not found")
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &body); err != nil {
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Name != "" {
		t.Name = body.Name
		if err := h.team.Update(r.Context(), t); err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	JSON(w, http.StatusOK, t)
}

func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	teamID := chi.URLParam(r, "teamID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	t, err := h.team.GetByID(r.Context(), teamID)
	if err != nil || t.OrgID != orgID {
		Error(w, http.StatusNotFound, "team not found")
		return
	}
	if err := h.team.Delete(r.Context(), teamID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
