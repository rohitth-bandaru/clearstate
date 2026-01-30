package model

import "time"

type User struct {
	ID        string     `db:"id" json:"id"`
	Email     string     `db:"email" json:"email"`
	Name      string     `db:"name" json:"name"`
	AvatarURL string     `db:"avatar_url" json:"avatar_url"`
	GoogleID  *string    `db:"google_id" json:"google_id,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type UserPublic struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}
