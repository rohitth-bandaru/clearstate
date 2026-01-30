package handler

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/clearstatus/backend/internal/model"
	"github.com/clearstatus/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type PublicHandler struct {
	org      *store.OrganizationStore
	service  *store.ServiceStore
	incident *store.IncidentStore
}

func NewPublicHandler(org *store.OrganizationStore, svc *store.ServiceStore, inc *store.IncidentStore) *PublicHandler {
	return &PublicHandler{org: org, service: svc, incident: inc}
}

// Status returns full status for public page (no auth). Used for polling.
func (h *PublicHandler) Status(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		Error(w, http.StatusBadRequest, "slug required")
		return
	}
	o, err := h.org.GetBySlug(r.Context(), slug)
	if err != nil {
		Error(w, http.StatusNotFound, "organization not found")
		return
	}
	services, err := h.service.ListByOrgID(r.Context(), o.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	incidents, err := h.incident.ListRecentByOrgID(r.Context(), o.ID, 30)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	active, _ := h.incident.ListActiveByOrgID(r.Context(), o.ID)

	resp := model.PublicStatusResponse{
		Org:       model.PublicOrg{ID: o.ID, Name: o.Name, Slug: o.Slug},
		Services:  make([]model.PublicService, 0, len(services)),
		Incidents: make([]model.PublicIncident, 0, len(incidents)),
		Timeline:  buildTimeline(r.Context(), h.incident, o.ID, services, incidents),
		Summary:   computeSummary(services, active),
	}
	for _, s := range services {
		resp.Services = append(resp.Services, model.PublicService{
			ID: s.ID, Name: s.Name, Slug: s.Slug, Description: s.Description, Status: s.Status,
		})
	}
	for _, inc := range incidents {
		serviceIDs, _ := h.incident.GetServiceIDs(r.Context(), inc.ID)
		updates, _ := h.incident.GetUpdates(r.Context(), inc.ID)
		pub := model.PublicIncident{
			ID:         inc.ID,
			Title:      inc.Title,
			Status:     inc.Status,
			Type:       inc.Type,
			CreatedAt:  inc.CreatedAt,
			ResolvedAt: inc.ResolvedAt,
			ServiceIDs: serviceIDs,
			Updates:    make([]model.PublicIncidentUpdate, 0, len(updates)),
		}
		for _, u := range updates {
			pub.Updates = append(pub.Updates, model.PublicIncidentUpdate{Message: u.Message, Status: u.Status, CreatedAt: u.CreatedAt})
		}
		resp.Incidents = append(resp.Incidents, pub)
	}
	JSON(w, http.StatusOK, resp)
}

func computeSummary(services []model.Service, active []model.Incident) string {
	if len(active) > 0 {
		return "Active incidents"
	}
	worst := "operational"
	for _, s := range services {
		switch s.Status {
		case model.ServiceStatusMajorOutage:
			worst = "Major Outage"
		case model.ServiceStatusPartialOutage:
			if worst != "Major Outage" {
				worst = "Partial Outage"
			}
		case model.ServiceStatusDegraded:
			if worst == "operational" {
				worst = "Degraded Performance"
			}
		}
	}
	if worst == "operational" {
		return "All Systems Operational"
	}
	return worst
}

func buildTimeline(ctx context.Context, incStore *store.IncidentStore, orgID string, services []model.Service, incidents []model.Incident) []model.PublicTimelineEntry {
	type entry struct {
		at   time.Time
		item model.PublicTimelineEntry
	}
	var entries []entry
	for _, inc := range incidents {
		entries = append(entries, entry{inc.CreatedAt, model.PublicTimelineEntry{
			Type:       "incident",
			IncidentID: &inc.ID,
			Title:      &inc.Title,
			CreatedAt:  inc.CreatedAt,
		}})
		updates, _ := incStore.GetUpdates(ctx, inc.ID)
		for _, u := range updates {
			msg := u.Message
			entries = append(entries, entry{u.CreatedAt, model.PublicTimelineEntry{
				Type:       "incident_update",
				IncidentID: &inc.ID,
				Message:    &msg,
				CreatedAt:  u.CreatedAt,
			}})
		}
		if inc.ResolvedAt != nil {
			entries = append(entries, entry{*inc.ResolvedAt, model.PublicTimelineEntry{
				Type:       "incident_resolved",
				IncidentID: &inc.ID,
				Title:      &inc.Title,
				CreatedAt:  *inc.ResolvedAt,
			}})
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].at.After(entries[j].at) })
	out := make([]model.PublicTimelineEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.item)
	}
	if len(out) > 50 {
		out = out[:50]
	}
	return out
}
