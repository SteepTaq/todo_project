-- +migrate Down
DROP TABLE IF EXISTS tasks;
DROP EXTENSION IF EXISTS "pgcrypto";