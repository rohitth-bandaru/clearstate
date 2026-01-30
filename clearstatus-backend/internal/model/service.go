package model

import "time"

const (
	ServiceStatusOperational   = "operational"
	ServiceStatusDegraded      = "degraded"
	ServiceStatusPartialOutage = "partial_outage"
	ServiceStatusMajorOutage   = "major_outage"
)

type Service struct {
	ID          string    `db:"id" json:"id"`
	OrgID       string    `db:"org_id" json:"org_id"`
	TeamID      string    `db:"team_id" json:"team_id"`
	Name        string    `db:"name" json:"name"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	SortOrder   int       `db:"sort_order" json:"sort_order"`
	Version     int       `db:"version" json:"version"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type ServiceCreate struct {
	OrgID       string
	TeamID      string
	Name        string
	Slug        string
	Description string
	Status      string
	SortOrder   int
}

type ServiceUpdate struct {
	TeamID      *string
	Name        *string
	Slug        *string
	Description *string
	Status      *string
	SortOrder   *int
	Version     int // required for optimistic locking
}
