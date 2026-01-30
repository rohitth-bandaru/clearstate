package model

import "time"

type Organization struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type OrganizationMember struct {
	OrgID     string    `db:"org_id" json:"org_id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type OrganizationMemberWithUser struct {
	OrgID     string    `db:"org_id" json:"org_id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	AvatarURL string    `db:"avatar_url" json:"avatar_url"`
}
