-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TYPE ticket_status AS ENUM (
  'init',
  'open',
  'closed'
);

CREATE TYPE complaint_info AS (
    sender_name VARCHAR(150),
    sender_phone VARCHAR(20),
    sender_email VARCHAR(150),
    geo_location    GEOGRAPHY(Point, 4326)
);

CREATE TABLE tickets (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Info
    status ticket_status NOT NULL DEFAULT 'init',
    complaints complaint_info[] NOT NULL DEFAULT '{}',
    description TEXT NOT NULL,
    is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
    subcategory_id INT NOT NULL REFERENCES subcategories(id) ON DELETE RESTRICT,
    department_id INT REFERENCES departments(id) ON DELETE SET NULL,

    -- Semantic embedding
    embedding   VECTOR(768) NOT NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- CREATE INDEX tickets_location_idx
-- ON tickets
-- USING GIST (location);

CREATE INDEX tickets_embedding_idx
ON tickets
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- +goose Down
DROP TABLE tickets;
DROP TYPE ticket_status;
DROP TYPE complaint_info;

DROP EXTENSION "uuid-ossp";
DROP EXTENSION pgcrypto;
DROP EXTENSION vector;
