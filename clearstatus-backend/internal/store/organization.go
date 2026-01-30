package store

import (
	"context"
	"github.com/clearstatus/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrganizationStore struct {
	db *sqlx.DB
}

func NewOrganizationStore(db *sqlx.DB) *OrganizationStore { return &OrganizationStore{db: db} }

func (s *OrganizationStore) Create(ctx context.Context, o *model.Organization) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	query := `INSERT INTO organizations (id, name, slug, created_at, updated_at) VALUES (?, ?, ?, NOW(3), NOW(3))`
	_, err := s.db.ExecContext(ctx, query, o.ID, o.Name, o.Slug)
	return err
}

func (s *OrganizationStore) GetByID(ctx context.Context, id string) (*model.Organization, error) {
	var o model.Organization
	err := s.db.GetContext(ctx, &o, `SELECT id, name, slug, created_at, updated_at FROM organizations WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *OrganizationStore) GetBySlug(ctx context.Context, slug string) (*model.Organization, error) {
	var o model.Organization
	err := s.db.GetContext(ctx, &o, `SELECT id, name, slug, created_at, updated_at FROM organizations WHERE slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *OrganizationStore) ListByUserID(ctx context.Context, userID string) ([]model.Organization, error) {
	var list []model.Organization
	err := s.db.SelectContext(ctx, &list,
		`SELECT o.id, o.name, o.slug, o.created_at, o.updated_at FROM organizations o
		 INNER JOIN organization_members m ON m.org_id = o.id WHERE m.user_id = ? ORDER BY o.name`,
		userID)
	return list, err
}

func (s *OrganizationStore) Update(ctx context.Context, o *model.Organization) error {
	_, err := s.db.ExecContext(ctx, `UPDATE organizations SET name = ?, slug = ?, updated_at = NOW(3) WHERE id = ?`, o.Name, o.Slug, o.ID)
	return err
}

func (s *OrganizationStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM organizations WHERE id = ?`, id)
	return err
}

func (s *OrganizationStore) AddMember(ctx context.Context, orgID, userID, role string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO organization_members (org_id, user_id, role, created_at) VALUES (?, ?, ?, NOW(3))`, orgID, userID, role)
	return err
}

func (s *OrganizationStore) RemoveMember(ctx context.Context, orgID, userID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM organization_members WHERE org_id = ? AND user_id = ?`, orgID, userID)
	return err
}

func (s *OrganizationStore) GetMemberRole(ctx context.Context, orgID, userID string) (string, error) {
	var role string
	err := s.db.GetContext(ctx, &role, `SELECT role FROM organization_members WHERE org_id = ? AND user_id = ?`, orgID, userID)
	return role, err
}

func (s *OrganizationStore) ListMembers(ctx context.Context, orgID string) ([]model.OrganizationMemberWithUser, error) {
	var list []model.OrganizationMemberWithUser
	err := s.db.SelectContext(ctx, &list,
		`SELECT m.org_id, m.user_id, m.role, m.created_at, u.email, u.name, u.avatar_url
		 FROM organization_members m INNER JOIN users u ON u.id = m.user_id WHERE m.org_id = ? ORDER BY m.created_at`,
		orgID)
	return list, err
}

func (s *OrganizationStore) UserInOrg(ctx context.Context, orgID, userID string) (bool, error) {
	var n int
	err := s.db.GetContext(ctx, &n, `SELECT 1 FROM organization_members WHERE org_id = ? AND user_id = ?`, orgID, userID)
	return n == 1, err
}
