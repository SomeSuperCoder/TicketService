-- +goose Up
CREATE TYPE user_role AS ENUM ('admin', 'org', 'executor');
CREATE TYPE user_status AS ENUM ('active', 'blocked');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    role user_role NOT NULL DEFAULT 'executor',
    status user_status NOT NULL DEFAULT 'active',
    department_id INT REFERENCES departments(id) ON DELETE SET NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    middle_name VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_department ON users(department_id);

-- +goose Down
DROP TABLE users;
DROP TYPE user_status;
DROP TYPE user_role;
