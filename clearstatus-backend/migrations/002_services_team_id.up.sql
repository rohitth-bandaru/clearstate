-- Add team_id to services so every service belongs to a team
ALTER TABLE services ADD COLUMN team_id CHAR(36) NULL AFTER org_id;

-- For each org that has services, ensure a Default team exists and assign services to it
INSERT INTO teams (id, org_id, name, created_at, updated_at)
SELECT UUID(), d.org_id, 'Default', NOW(3), NOW(3)
FROM (SELECT DISTINCT org_id AS org_id FROM services) d;

UPDATE services s
SET s.team_id = (SELECT t.id FROM teams t WHERE t.org_id = s.org_id ORDER BY t.created_at DESC LIMIT 1);

ALTER TABLE services MODIFY team_id CHAR(36) NOT NULL;
ALTER TABLE services ADD CONSTRAINT fk_services_team FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE RESTRICT;
CREATE INDEX idx_services_team ON services(team_id);
