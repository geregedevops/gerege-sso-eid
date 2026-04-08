CREATE TABLE IF NOT EXISTS dan_clients (
    id            TEXT PRIMARY KEY,
    secret_hash   TEXT NOT NULL,
    name          TEXT NOT NULL,
    callback_urls TEXT[] NOT NULL DEFAULT '{}',
    active        BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_dan_clients_active ON dan_clients(active);
