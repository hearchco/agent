-- +goose Up
-- general results
CREATE TABLE results (
  id BIGSERIAL PRIMARY KEY,
  query TEXT NOT NULL,
  url TEXT NOT NULL,
  rank BIGINT NOT NULL,
  score DOUBLE PRECISION NOT NULL,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- image part of results
CREATE TABLE image_results (
  id BIGSERIAL PRIMARY KEY,
  image_original_height BIGINT NOT NULL,
  image_original_width BIGINT NOT NULL,
  image_thumbnail_height BIGINT NOT NULL,
  image_thumbnail_width BIGINT NOT NULL,
  image_thumbnail_url TEXT NOT NULL,
  image_source TEXT NOT NULL,
  image_source_url TEXT NOT NULL,
  result_id BIGINT NOT NULL,
  FOREIGN KEY (result_id) REFERENCES results (id) ON DELETE CASCADE
);

-- engine ranks
CREATE TABLE engine_ranks (
  id BIGSERIAL PRIMARY KEY,
  engine_name TEXT NOT NULL,
  engine_rank BIGINT NOT NULL,
  engine_page BIGINT NOT NULL,
  engine_on_page_rank BIGINT NOT NULL,
  result_id BIGINT NOT NULL,
  FOREIGN KEY (result_id) REFERENCES results (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE engine_ranks;
DROP TABLE image_results;
DROP TABLE results;