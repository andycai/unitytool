-- +migrate Up
ALTER TABLE tasks ADD COLUMN enable_cron TINYINT DEFAULT 0;
ALTER TABLE tasks ADD COLUMN cron_expr VARCHAR(100);

-- +migrate Down
ALTER TABLE tasks DROP COLUMN enable_cron;
ALTER TABLE tasks DROP COLUMN cron_expr; 