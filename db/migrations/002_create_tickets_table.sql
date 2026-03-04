-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TYPE ticket_status AS ENUM (
  'none',
  'init',
  'open',
  'closed'
);

CREATE TABLE tickets (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    status ticket_status NOT NULL DEFAULT 'init',
    description TEXT NOT NULL,
    subcategory_id INT NOT NULL REFERENCES subcategories(id) ON DELETE RESTRICT,
    department_id INT REFERENCES departments(id) ON DELETE SET NULL,

    -- Semantic embedding
    embedding   VECTOR(768) NOT NULL,

    is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE complaint_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    ticket UUID NOT NULL REFERENCES tickets (id) ON DELETE CASCADE,
    description TEXT NOT NULL,

    sender_name VARCHAR(150) NOT NULL,
    sender_phone VARCHAR(20),
    sender_email VARCHAR(150),
    geo_location GEOGRAPHY(Point, 4326),
    
    CONSTRAINT chk_sender_name_not_empty CHECK (LENGTH(TRIM(sender_name)) > 0),
    CONSTRAINT chk_sender_phone_format CHECK (
        sender_phone IS NULL OR 
        sender_phone ~ '^[0-9+\-\s()]{10,20}$'
    ),
    CONSTRAINT chk_sender_email_format CHECK (
        sender_email IS NULL OR 
        sender_email ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    ),
    CONSTRAINT chk_at_least_one_contact CHECK (
        sender_phone IS NOT NULL OR 
        sender_email IS NOT NULL
    )
);

CREATE TABLE ticket_tags (
  ticket UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
  tag INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  CONSTRAINT uq_ticket_tag UNIQUE (ticket, tag)
);

CREATE INDEX complaint_details_location_idx
ON complaint_details
USING GIST (geo_location);

CREATE INDEX tickets_embedding_idx
ON tickets
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- +goose Down
DROP TABLE complaint_details ;
DROP TABLE tickets;
DROP TYPE ticket_status;
DROP TABLE ticket_tags;

DROP EXTENSION "uuid-ossp";
DROP EXTENSION pgcrypto;
DROP EXTENSION vector;
