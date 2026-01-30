package store

import (
	"context"
	"github.com/clearstatus/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore { return &UserStore{db: db} }

func (s *UserStore) Create(ctx context.Context, u *model.User) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	query := `INSERT INTO users (id, email, name, avatar_url, google_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(3), NOW(3))`
	_, err := s.db.ExecContext(ctx, query, u.ID, u.Email, u.Name, u.AvatarURL, u.GoogleID)
	return err
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	err := s.db.GetContext(ctx, &u, `SELECT id, email, name, avatar_url, google_id, created_at, updated_at FROM users WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := s.db.GetContext(ctx, &u, `SELECT id, email, name, avatar_url, google_id, created_at, updated_at FROM users WHERE email = ?`, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserStore) GetByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	var u model.User
	err := s.db.GetContext(ctx, &u, `SELECT id, email, name, avatar_url, google_id, created_at, updated_at FROM users WHERE google_id = ?`, googleID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserStore) Update(ctx context.Context, u *model.User) error {
	query := `UPDATE users SET name = ?, avatar_url = ?, google_id = ?, updated_at = NOW(3) WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, u.Name, u.AvatarURL, u.GoogleID, u.ID)
	return err
}
