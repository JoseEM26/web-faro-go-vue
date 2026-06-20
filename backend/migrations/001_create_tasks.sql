-- Migracion 001: Crear tabla tasks
-- GORM aplica esto automaticamente con AutoMigrate al iniciar el servidor.
-- Este archivo documenta el esquema para control de versiones.

-- UP
CREATE TABLE IF NOT EXISTS tasks (
    id          BIGSERIAL    PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         DEFAULT '',
    completed   BOOLEAN      NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks(completed);

-- DOWN (rollback manual si es necesario)
-- DROP TABLE IF EXISTS tasks;
