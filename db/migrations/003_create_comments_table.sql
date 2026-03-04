-- +goose Up
CREATE TABLE comments (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    ticket UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    message TEXT NOT NULL
);

-- +goose Down
DROP TABLE comments;
