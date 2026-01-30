package store

import (
	"context"
	"github.com/clearstatus/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IncidentStore struct {
	db *sqlx.DB
}

func NewIncidentStore(db *sqlx.DB) *IncidentStore { return &IncidentStore{db: db} }

func (s *IncidentStore) Create(ctx context.Context, in model.IncidentCreate) (*model.Incident, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	id := uuid.New().String()
	_, err = tx.ExecContext(ctx, `INSERT INTO incidents (id, org_id, title, status, type, version, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 1, NOW(3), NOW(3))`,
		id, in.OrgID, in.Title, in.Status, in.Type)
	if err != nil {
		return nil, err
	}
	for _, sid := range in.ServiceIDs {
		_, err = tx.ExecContext(ctx, `INSERT INTO incident_services (incident_id, service_id) VALUES (?, ?)`, id, sid)
		if err != nil {
			return nil, err
		}
	}
	if in.Message != "" {
		upID := uuid.New().String()
		_, err = tx.ExecContext(ctx, `INSERT INTO incident_updates (id, incident_id, message, status, created_at) VALUES (?, ?, ?, ?, NOW(3))`, upID, id, in.Message, in.Status)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *IncidentStore) GetByID(ctx context.Context, id string) (*model.Incident, error) {
	var inc model.Incident
	err := s.db.GetContext(ctx, &inc, `SELECT id, org_id, title, status, type, version, created_at, updated_at, resolved_at FROM incidents WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &inc, nil
}

func (s *IncidentStore) ListByOrgID(ctx context.Context, orgID string, limit int) ([]model.Incident, error) {
	if limit <= 0 {
		limit = 50
	}
	var list []model.Incident
	err := s.db.SelectContext(ctx, &list, `SELECT id, org_id, title, status, type, version, created_at, updated_at, resolved_at FROM incidents WHERE org_id = ? ORDER BY created_at DESC LIMIT ?`, orgID, limit)
	return list, err
}

func (s *IncidentStore) GetServiceIDs(ctx context.Context, incidentID string) ([]string, error) {
	var ids []string
	err := s.db.SelectContext(ctx, &ids, `SELECT service_id FROM incident_services WHERE incident_id = ?`, incidentID)
	return ids, err
}

func (s *IncidentStore) GetUpdates(ctx context.Context, incidentID string) ([]model.IncidentUpdateMessage, error) {
	var list []model.IncidentUpdateMessage
	err := s.db.SelectContext(ctx, &list, `SELECT id, incident_id, message, status, created_at FROM incident_updates WHERE incident_id = ? ORDER BY created_at`, incidentID)
	return list, err
}

func (s *IncidentStore) Update(ctx context.Context, id string, in model.IncidentUpdate) (*model.Incident, error) {
	query := `UPDATE incidents SET updated_at = NOW(3), version = version + 1`
	args := []interface{}{}
	if in.Title != nil {
		query += `, title = ?`
		args = append(args, *in.Title)
	}
	if in.Status != nil {
		query += `, status = ?`
		args = append(args, *in.Status)
		if *in.Status == model.IncidentStatusResolved {
			query += `, resolved_at = NOW(3)`
		}
	}
	query += ` WHERE id = ? AND version = ?`
	args = append(args, id, in.Version)
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, ErrConflict
	}
	return s.GetByID(ctx, id)
}

func (s *IncidentStore) AddUpdate(ctx context.Context, incidentID, message, status string) (*model.IncidentUpdateMessage, error) {
	upID := uuid.New().String()
	_, err := s.db.ExecContext(ctx, `INSERT INTO incident_updates (id, incident_id, message, status, created_at) VALUES (?, ?, ?, ?, NOW(3))`, upID, incidentID, message, status)
	if err != nil {
		return nil, err
	}
	var u model.IncidentUpdateMessage
	err = s.db.GetContext(ctx, &u, `SELECT id, incident_id, message, status, created_at FROM incident_updates WHERE id = ?`, upID)
	return &u, err
}

func (s *IncidentStore) SetIncidentServices(ctx context.Context, incidentID string, serviceIDs []string) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `DELETE FROM incident_services WHERE incident_id = ?`, incidentID)
	if err != nil {
		return err
	}
	for _, sid := range serviceIDs {
		_, err = tx.ExecContext(ctx, `INSERT INTO incident_services (incident_id, service_id) VALUES (?, ?)`, incidentID, sid)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *IncidentStore) ListActiveByOrgID(ctx context.Context, orgID string) ([]model.Incident, error) {
	var list []model.Incident
	err := s.db.SelectContext(ctx, &list, `SELECT id, org_id, title, status, type, version, created_at, updated_at, resolved_at FROM incidents WHERE org_id = ? AND status != 'resolved' ORDER BY created_at DESC`, orgID)
	return list, err
}

func (s *IncidentStore) ListRecentByOrgID(ctx context.Context, orgID string, limit int) ([]model.Incident, error) {
	if limit <= 0 {
		limit = 30
	}
	var list []model.Incident
	err := s.db.SelectContext(ctx, &list, `SELECT id, org_id, title, status, type, version, created_at, updated_at, resolved_at FROM incidents WHERE org_id = ? ORDER BY created_at DESC LIMIT ?`, orgID, limit)
	return list, err
}

func (s *IncidentStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM incidents WHERE id = ?`, id)
	return err
}

