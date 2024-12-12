-- +migrate Up
ALTER TABLE users ADD COLUMN has_changed_pwd BOOLEAN DEFAULT false;

-- +migrate Down
ALTER TABLE users DROP COLUMN has_changed_pwd; 