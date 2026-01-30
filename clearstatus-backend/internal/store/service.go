package store

import (
	"context"
	"github.com/clearstatus/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ServiceStore struct {
	db *sqlx.DB
}

func NewServiceStore(db *sqlx.DB) *ServiceStore { return &ServiceStore{db: db} }

func (s *ServiceStore) Create(ctx context.Context, in model.ServiceCreate) (*model.Service, error) {
	id := uuid.New().String()
	query := `INSERT INTO services (id, org_id, team_id, name, slug, description, status, sort_order, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, NOW(3), NOW(3))`
	_, err := s.db.ExecContext(ctx, query, id, in.OrgID, in.TeamID, in.Name, in.Slug, in.Description, in.Status, in.SortOrder)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *ServiceStore) GetByID(ctx context.Context, id string) (*model.Service, error) {
	var svc model.Service
	err := s.db.GetContext(ctx, &svc, `SELECT id, org_id, team_id, name, slug, description, status, sort_order, version, created_at, updated_at FROM services WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &svc, nil
}

func (s *ServiceStore) GetByOrgAndSlug(ctx context.Context, orgID, slug string) (*model.Service, error) {
	var svc model.Service
	err := s.db.GetContext(ctx, &svc, `SELECT id, org_id, team_id, name, slug, description, status, sort_order, version, created_at, updated_at FROM services WHERE org_id = ? AND slug = ?`, orgID, slug)
	if err != nil {
		return nil, err
	}
	return &svc, nil
}

func (s *ServiceStore) ListByOrgID(ctx context.Context, orgID string) ([]model.Service, error) {
	var list []model.Service
	err := s.db.SelectContext(ctx, &list, `SELECT id, org_id, team_id, name, slug, description, status, sort_order, version, created_at, updated_at FROM services WHERE org_id = ? ORDER BY sort_order, name`, orgID)
	return list, err
}

func (s *ServiceStore) Update(ctx context.Context, id string, in model.ServiceUpdate) (*model.Service, error) {
	// Optimistic locking: WHERE id = ? AND version = ?
	query := `UPDATE services SET updated_at = NOW(3), version = version + 1`
	args := []interface{}{}
	if in.TeamID != nil {
		query += `, team_id = ?`
		args = append(args, *in.TeamID)
	}
	if in.Name != nil {
		query += `, name = ?`
		args = append(args, *in.Name)
	}
	if in.Slug != nil {
		query += `, slug = ?`
		args = append(args, *in.Slug)
	}
	if in.Description != nil {
		query += `, description = ?`
		args = append(args, *in.Description)
	}
	if in.Status != nil {
		query += `, status = ?`
		args = append(args, *in.Status)
	}
	if in.SortOrder != nil {
		query += `, sort_order = ?`
		args = append(args, *in.SortOrder)
	}
	query += ` WHERE id = ? AND version = ?`
	args = append(args, id, in.Version)
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, ErrConflict // 409
	}
	return s.GetByID(ctx, id)
}

func (s *ServiceStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM services WHERE id = ?`, id)
	return err
}

func (s *ServiceStore) RecordStatusHistory(ctx context.Context, serviceID, status string) error {
	id := uuid.New().String()
	_, err := s.db.ExecContext(ctx, `INSERT INTO status_history (id, service_id, status, created_at) VALUES (?, ?, ?, NOW(3))`, id, serviceID, status)
	return err
}
