-- +goose Up

CREATE TABLE feeds (
  id uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR NOT NULL,
  url VARCHAR UNIQUE NOT NULL,
  user_id uuid NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;