ALTER TABLE services DROP FOREIGN KEY fk_services_team;
DROP INDEX idx_services_team ON services;
ALTER TABLE services DROP COLUMN team_id;
