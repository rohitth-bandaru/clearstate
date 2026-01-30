package model

import "time"

const (
	IncidentStatusInvestigating = "investigating"
	IncidentStatusIdentified    = "identified"
	IncidentStatusMonitoring    = "monitoring"
	IncidentStatusResolved       = "resolved"
	IncidentTypeIncident        = "incident"
	IncidentTypeMaintenance     = "maintenance"
)

type Incident struct {
	ID         string     `db:"id" json:"id"`
	OrgID      string     `db:"org_id" json:"org_id"`
	Title      string     `db:"title" json:"title"`
	Status     string     `db:"status" json:"status"`
	Type       string     `db:"type" json:"type"`
	Version    int        `db:"version" json:"version"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	ResolvedAt *time.Time `db:"resolved_at" json:"resolved_at,omitempty"`
}

type IncidentCreate struct {
	OrgID     string
	Title     string
	Status    string
	Type      string
	ServiceIDs []string
	Message   string
}

type IncidentUpdate struct {
	Title     *string
	Status    *string
	Version   int // required for optimistic locking
}

type IncidentUpdateMessage struct {
	ID        string    `db:"id" json:"id"`
	IncidentID string   `db:"incident_id" json:"incident_id"`
	Message   string    `db:"message" json:"message"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type IncidentWithServices struct {
	Incident
	ServiceIDs []string `json:"service_ids"`
	Updates    []IncidentUpdateMessage `json:"updates,omitempty"`
}
