-- +goose Up
CREATE TYPE history_action AS ENUM (
  'created',
  'status_changed',
  'department_changed',
  'tags_added',
  'tags_removed',
  'comment_added',
  'merged',
  'deleted'
);

CREATE TABLE ticket_history (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id   UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    action      history_action NOT NULL,
    
    -- Store the changes as JSONB for flexibility
    old_value   JSONB,
    new_value   JSONB,
    
    -- User who made the change
    user_id     UUID,
    user_name   VARCHAR(150),
    user_email  VARCHAR(150),
    
    -- Optional description
    description TEXT,
    
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ticket_history_ticket_id_idx ON ticket_history(ticket_id);
CREATE INDEX ticket_history_created_at_idx ON ticket_history(created_at DESC);

-- +goose Down
DROP TABLE ticket_history;
DROP TYPE history_action;
