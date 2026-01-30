package handler

import (
	"net/http"

	"github.com/clearstatus/backend/internal/middleware"
	"github.com/clearstatus/backend/internal/model"
	"github.com/clearstatus/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type IncidentHandler struct {
	org      *store.OrganizationStore
	incident *store.IncidentStore
}

func NewIncidentHandler(org *store.OrganizationStore, inc *store.IncidentStore) *IncidentHandler {
	return &IncidentHandler{org: org, incident: inc}
}

func (h *IncidentHandler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{incidentID}", h.Get)
	r.Put("/{incidentID}", h.Update)
	r.Delete("/{incidentID}", h.Delete)
	r.Post("/{incidentID}/updates", h.AddUpdate)
}

func (h *IncidentHandler) ensureOrgMember(w http.ResponseWriter, r *http.Request, orgID string) bool {
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

func (h *IncidentHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	list, err := h.incident.ListByOrgID(r.Context(), orgID, 50)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, list)
}

func (h *IncidentHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	var body struct {
		Title      string   `json:"title"`
		Status     string   `json:"status"`
		Type       string   `json:"type"`
		ServiceIDs []string `json:"service_ids"`
		Message    string   `json:"message"`
	}
	if err := decodeJSON(r, &body); err != nil {
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Title == "" {
		Error(w, http.StatusBadRequest, "title required")
		return
	}
	if body.Status == "" {
		body.Status = model.IncidentStatusInvestigating
	}
	if body.Type == "" {
		body.Type = model.IncidentTypeIncident
	}
	if body.Type != model.IncidentTypeIncident && body.Type != model.IncidentTypeMaintenance {
		Error(w, http.StatusBadRequest, "type must be incident or maintenance")
		return
	}
	if !validIncidentStatus(body.Status) {
		Error(w, http.StatusBadRequest, "invalid status")
		return
	}
	in := model.IncidentCreate{OrgID: orgID, Title: body.Title, Status: body.Status, Type: body.Type, ServiceIDs: body.ServiceIDs, Message: body.Message}
	inc, err := h.incident.Create(r.Context(), in)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, inc)
}

func (h *IncidentHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	incidentID := chi.URLParam(r, "incidentID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	inc, err := h.incident.GetByID(r.Context(), incidentID)
	if err != nil {
		Error(w, http.StatusNotFound, "incident not found")
		return
	}
	if inc.OrgID != orgID {
		Error(w, http.StatusNotFound, "incident not found")
		return
	}
	serviceIDs, _ := h.incident.GetServiceIDs(r.Context(), incidentID)
	updates, _ := h.incident.GetUpdates(r.Context(), incidentID)
	out := model.IncidentWithServices{Incident: *inc, ServiceIDs: serviceIDs, Updates: updates}
	JSON(w, http.StatusOK, out)
}

func (h *IncidentHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	incidentID := chi.URLParam(r, "incidentID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	inc, err := h.incident.GetByID(r.Context(), incidentID)
	if err != nil || inc.OrgID != orgID {
		Error(w, http.StatusNotFound, "incident not found")
		return
	}
	var body struct {
		Title   *string  `json:"title"`
		Status  *string  `json:"status"`
		Version int      `json:"version"`
	}
	if err := decodeJSON(r, &body); err != nil {
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Version == 0 {
		body.Version = inc.Version
	}
	in := model.IncidentUpdate{Version: body.Version}
	if body.Title != nil {
		in.Title = body.Title
	}
	if body.Status != nil {
		if !validIncidentStatus(*body.Status) {
			Error(w, http.StatusBadRequest, "invalid status")
			return
		}
		in.Status = body.Status
	}
	updated, err := h.incident.Update(r.Context(), incidentID, in)
	if err != nil {
		if err == store.ErrConflict {
			Error(w, http.StatusConflict, "resource was modified; refetch and retry")
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, updated)
}

func (h *IncidentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	incidentID := chi.URLParam(r, "incidentID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	inc, err := h.incident.GetByID(r.Context(), incidentID)
	if err != nil || inc.OrgID != orgID {
		Error(w, http.StatusNotFound, "incident not found")
		return
	}
	if err := h.incident.Delete(r.Context(), incidentID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *IncidentHandler) AddUpdate(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	incidentID := chi.URLParam(r, "incidentID")
	if !h.ensureOrgMember(w, r, orgID) {
		return
	}
	inc, err := h.incident.GetByID(r.Context(), incidentID)
	if err != nil || inc.OrgID != orgID {
		Error(w, http.StatusNotFound, "incident not found")
		return
	}
	var body struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	if err := decodeJSON(r, &body); err != nil {
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Message == "" {
		Error(w, http.StatusBadRequest, "message required")
		return
	}
	if body.Status == "" {
		body.Status = inc.Status
	}
	if !validIncidentStatus(body.Status) {
		Error(w, http.StatusBadRequest, "invalid status")
		return
	}
	up, err := h.incident.AddUpdate(r.Context(), incidentID, body.Message, body.Status)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, up)
}

func validIncidentStatus(s string) bool {
	return s == model.IncidentStatusInvestigating || s == model.IncidentStatusIdentified || s == model.IncidentStatusMonitoring || s == model.IncidentStatusResolved
}
