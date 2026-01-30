package model

import "time"

// PublicStatusResponse is returned by GET /api/public/orgs/:slug/status
type PublicStatusResponse struct {
	Org       PublicOrg        `json:"org"`
	Services  []PublicService  `json:"services"`
	Incidents []PublicIncident `json:"incidents"`
	Timeline  []PublicTimelineEntry `json:"timeline"`
	Summary   string           `json:"summary"` // e.g. "All Systems Operational"
}

type PublicOrg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type PublicService struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type PublicIncident struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	ServiceIDs []string  `json:"service_ids"`
	Updates    []PublicIncidentUpdate `json:"updates"`
}

type PublicIncidentUpdate struct {
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type PublicTimelineEntry struct {
	Type      string    `json:"type"` // "service_status" | "incident"
	ServiceID *string   `json:"service_id,omitempty"`
	ServiceName *string `json:"service_name,omitempty"`
	Status    *string  `json:"status,omitempty"`
	IncidentID *string `json:"incident_id,omitempty"`
	Title     *string  `json:"title,omitempty"`
	Message   *string  `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
