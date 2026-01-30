package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/clearstatus/backend/internal/middleware"
	"github.com/clearstatus/backend/internal/model"
	"github.com/clearstatus/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

var serviceSlugRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)

type ServiceHandler struct {
	org     *store.OrganizationStore
	team    *store.TeamStore
	service *store.ServiceStore
}

func NewServiceHandler(org *store.OrganizationStore, team *store.TeamStore, svc *store.ServiceStore) *ServiceHandler {
	return &ServiceHandler{org: org, team: team, service: svc}
}

func (h *ServiceHandler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{serviceID}", h.Get)
	r.Put("/{serviceID}", h.Update)
	r.Delete("/{serviceID}", h.Delete)
}

func (h *ServiceHandler) ensureOrgMember(w http.ResponseWriter, r *http.Request, orgID string) bool {
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

func (h *ServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	list, err := h.service.ListByOrgID(r.Context(), orgID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, list)
}

func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	var body struct {
		TeamID      string `json:"team_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := decodeJSON(r, &body); err != nil {
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.TeamID == "" {
		Error(w, http.StatusBadRequest, "team_id required")
		return
	}
	team, err := h.team.GetByID(r.Context(), body.TeamID)
	if err != nil || team.OrgID != orgID {
		Error(w, http.StatusBadRequest, "team not found or does not belong to this organization")
		return
	}
	if body.Name == "" {
		Error(w, http.StatusBadRequest, "name required")
		return
	}
	slug := slugify(body.Name)
	if slug == "" {
		slug = "service"
	}
	if !serviceSlugRe.MatchString(slug) {
		slug = "service"
	}
	if body.Status == "" {
		body.Status = model.ServiceStatusOperational
	}
	if !validServiceStatus(body.Status) {
		Error(w, http.StatusBadRequest, "invalid status")
		return
	}
	in := model.ServiceCreate{OrgID: orgID, TeamID: body.TeamID, Name: body.Name, Slug: slug, Description: body.Description, Status: body.Status, SortOrder: body.SortOrder}
	svc, err := h.service.Create(r.Context(), in)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			Error(w, http.StatusConflict, "slug already exists")
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = h.service.RecordStatusHistory(r.Context(), svc.ID, svc.Status)
	JSON(w, http.StatusCreated, svc)
}

func (h *ServiceHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	serviceID := chi.URLParam(r, "serviceID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	svc, err := h.service.GetByID(r.Context(), serviceID)
	if err != nil {
		Error(w, http.StatusNotFound, "service not found")
		return
	}
	if svc.OrgID != orgID {
		Error(w, http.StatusNotFound, "service not found")
		return
	}
	JSON(w, http.StatusOK, svc)
}

func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	serviceID := chi.URLParam(r, "serviceID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	svc, err := h.service.GetByID(r.Context(), serviceID)
	if err != nil || svc.OrgID != orgID {
		Error(w, http.StatusNotFound, "service not found")
		return
	}
	var body struct {
		TeamID      *string `json:"team_id"`
		Name        *string `json:"name"`
		Slug        *string `json:"slug"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
		SortOrder   *int    `json:"sort_order"`
		Version     int     `json:"version"`
	}
	if err := decodeJSON(r, &body); err != nil {
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Version == 0 {
		body.Version = svc.Version
	}
	in := model.ServiceUpdate{Version: body.Version}
	if body.TeamID != nil {
		team, err := h.team.GetByID(r.Context(), *body.TeamID)
		if err != nil || team.OrgID != orgID {
			Error(w, http.StatusBadRequest, "team not found or does not belong to this organization")
			return
		}
		in.TeamID = body.TeamID
	}
	if body.Name != nil {
		in.Name = body.Name
	}
	if body.Slug != nil {
		s := strings.TrimSpace(strings.ToLower(*body.Slug))
		in.Slug = &s
		if !serviceSlugRe.MatchString(s) {
			Error(w, http.StatusBadRequest, "invalid slug")
			return
		}
	}
	if body.Description != nil {
		in.Description = body.Description
	}
	if body.Status != nil {
		if !validServiceStatus(*body.Status) {
			Error(w, http.StatusBadRequest, "invalid status")
			return
		}
		in.Status = body.Status
	}
	if body.SortOrder != nil {
		in.SortOrder = body.SortOrder
	}
	updated, err := h.service.Update(r.Context(), serviceID, in)
	if err != nil {
		if err == store.ErrConflict {
			Error(w, http.StatusConflict, "resource was modified; refetch and retry")
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if body.Status != nil && *body.Status != svc.Status {
		_ = h.service.RecordStatusHistory(r.Context(), serviceID, *body.Status)
	}
	JSON(w, http.StatusOK, updated)
}

func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	serviceID := chi.URLParam(r, "serviceID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	svc, err := h.service.GetByID(r.Context(), serviceID)
	if err != nil || svc.OrgID != orgID {
		Error(w, http.StatusNotFound, "service not found")
		return
	}
	if err := h.service.Delete(r.Context(), serviceID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validServiceStatus(s string) bool {
	return s == model.ServiceStatusOperational || s == model.ServiceStatusDegraded || s == model.ServiceStatusPartialOutage || s == model.ServiceStatusMajorOutage
}
