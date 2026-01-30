package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/clearstatus/backend/internal/middleware"
	"github.com/clearstatus/backend/internal/model"
	"github.com/clearstatus/backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var slugRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)

type OrganizationHandler struct {
	org *store.OrganizationStore
}

func NewOrganizationHandler(org *store.OrganizationStore) *OrganizationHandler {
	return &OrganizationHandler{org: org}
}

func (h *OrganizationHandler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{orgID}", h.Get)
	r.Put("/{orgID}", h.Update)
	r.Delete("/{orgID}", h.Delete)
	r.Get("/{orgID}/members", h.ListMembers)
	r.Post("/{orgID}/members", h.AddMember)
	r.Delete("/{orgID}/members/{userID}", h.RemoveMember)
}

func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	u := middleware.UserFromContext(r.Context())
	if u == nil {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := h.org.ListByUserID(r.Context(), u.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, list)
}

func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if err := decodeJSON(r, &body); err != nil {
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	body.Slug = strings.TrimSpace(strings.ToLower(body.Slug))
	if body.Slug == "" {
		body.Slug = slugify(body.Name)
	}
	if !slugRe.MatchString(body.Slug) {
		Error(w, http.StatusBadRequest, "slug must be lowercase alphanumeric and hyphens")
		return
	}
	o := &model.Organization{ID: uuid.New().String(), Name: body.Name, Slug: body.Slug}
	if err := h.org.Create(r.Context(), o); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			Error(w, http.StatusConflict, "slug already exists")
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Add creator as owner
	_ = h.org.AddMember(r.Context(), o.ID, user.ID, "owner")
	JSON(w, http.StatusCreated, o)
}

func (h *OrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ok, err := h.org.UserInOrg(r.Context(), orgID, user.ID)
	if err != nil || !ok {
		Error(w, http.StatusNotFound, "organization not found")
		return
	}
	o, err := h.org.GetByID(r.Context(), orgID)
	if err != nil {
		Error(w, http.StatusNotFound, "organization not found")
		return
	}
	JSON(w, http.StatusOK, o)
}

func (h *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ok, _ := h.org.UserInOrg(r.Context(), orgID, user.ID)
	if !ok {
		Error(w, http.StatusNotFound, "organization not found")
		return
	}
	var body struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if err := decodeJSON(r, &body); err != nil {
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	o, err := h.org.GetByID(r.Context(), orgID)
	if err != nil {
		Error(w, http.StatusNotFound, "not found")
		return
	}
	if body.Name != "" {
		o.Name = body.Name
	}
	if body.Slug != "" {
		o.Slug = strings.TrimSpace(strings.ToLower(body.Slug))
		if !slugRe.MatchString(o.Slug) {
			Error(w, http.StatusBadRequest, "invalid slug")
			return
		}
	}
	if err := h.org.Update(r.Context(), o); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			Error(w, http.StatusConflict, "slug already exists")
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, o)
}

func (h *OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	role, err := h.org.GetMemberRole(r.Context(), orgID, user.ID)
	if err != nil || role != "owner" {
		Error(w, http.StatusForbidden, "only owner can delete")
		return
	}
	if err := h.org.Delete(r.Context(), orgID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrganizationHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ok, _ := h.org.UserInOrg(r.Context(), orgID, user.ID)
	if !ok {
		Error(w, http.StatusNotFound, "not found")
		return
	}
	list, err := h.org.ListMembers(r.Context(), orgID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, list)
}

func (h *OrganizationHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ok, _ := h.org.UserInOrg(r.Context(), orgID, user.ID)
	if !ok {
		Error(w, http.StatusNotFound, "not found")
		return
	}
	var body struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := decodeJSON(r, &body); err != nil || body.UserID == "" {
		Error(w, http.StatusBadRequest, "user_id required")
		return
	}
	if body.Role == "" {
		body.Role = "member"
	}
	if body.Role != "owner" && body.Role != "admin" && body.Role != "member" {
		Error(w, http.StatusBadRequest, "role must be owner, admin, or member")
		return
	}
	if err := h.org.AddMember(r.Context(), orgID, body.UserID, body.Role); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *OrganizationHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	userID := chi.URLParam(r, "userID")
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ok, _ := h.org.UserInOrg(r.Context(), orgID, user.ID)
	if !ok {
		Error(w, http.StatusNotFound, "not found")
		return
	}
	if err := h.org.RemoveMember(r.Context(), orgID, userID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' {
			if b.Len() > 0 && b.String()[b.Len()-1] != '-' {
				b.WriteRune('-')
			}
		}
	}
	return strings.Trim(b.String(), "-")
}
