package store

import (
	"context"
	"github.com/clearstatus/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeamStore struct {
	db *sqlx.DB
}

func NewTeamStore(db *sqlx.DB) *TeamStore { return &TeamStore{db: db} }

func (s *TeamStore) Create(ctx context.Context, t *model.Team) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO teams (id, org_id, name, created_at, updated_at) VALUES (?, ?, ?, NOW(3), NOW(3))`, t.ID, t.OrgID, t.Name)
	return err
}

func (s *TeamStore) GetByID(ctx context.Context, id string) (*model.Team, error) {
	var t model.Team
	err := s.db.GetContext(ctx, &t, `SELECT id, org_id, name, created_at, updated_at FROM teams WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TeamStore) ListByOrgID(ctx context.Context, orgID string) ([]model.Team, error) {
	var list []model.Team
	err := s.db.SelectContext(ctx, &list, `SELECT id, org_id, name, created_at, updated_at FROM teams WHERE org_id = ? ORDER BY name`, orgID)
	return list, err
}

func (s *TeamStore) Update(ctx context.Context, t *model.Team) error {
	_, err := s.db.ExecContext(ctx, `UPDATE teams SET name = ?, updated_at = NOW(3) WHERE id = ?`, t.Name, t.ID)
	return err
}

func (s *TeamStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM teams WHERE id = ?`, id)
	return err
}

func (s *TeamStore) AddMember(ctx context.Context, teamID, userID, role string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO team_members (team_id, user_id, role, created_at) VALUES (?, ?, ?, NOW(3))`, teamID, userID, role)
	return err
}

func (s *TeamStore) RemoveMember(ctx context.Context, teamID, userID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM team_members WHERE team_id = ? AND user_id = ?`, teamID, userID)
	return err
}
