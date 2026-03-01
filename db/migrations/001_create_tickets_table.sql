-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE tickets (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       TEXT NOT NULL,
    description TEXT NOT NULL,

    -- Geo location (lon/lat)
    location    GEOGRAPHY(Point, 4326) NOT NULL,

    -- Semantic embedding
    embedding   VECTOR(768) NOT NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX tickets_location_idx
ON tickets
USING GIST (location);

CREATE INDEX tickets_embedding_idx
ON tickets
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- +goose Down
DROP TABLE tickets;
DROP EXTENSION "uuid-ossp";
DROP EXTENSION pgcrypto;
DROP EXTENSION vector;
