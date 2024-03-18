-- +goose Up
CREATE TABLE posts (
  id uuid PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title VARCHAR NOT NULL,
  url VARCHAR UNIQUE NOT NULL,
  description TEXT NOT NULL,
  published_at TIMESTAMP NOT NULL,
  feed_id uuid NOT NULL,
  FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE
);


-- +goose Down

DROP TABLE posts;

