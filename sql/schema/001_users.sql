-- postgresql

-- +goose Up
CREATE TABLE users (
  id uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR NOT NULL
);

-- +goose Down
DROP TABLE users;
