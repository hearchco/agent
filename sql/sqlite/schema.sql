-- general results
CREATE TABLE results (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  query TEXT NOT NULL,
  url TEXT NOT NULL,
  rank INTEGER NOT NULL,
  score REAL NOT NULL,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- image part of results
CREATE TABLE image_results (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  image_original_height INTEGER NOT NULL,
  image_original_width INTEGER NOT NULL,
  image_thumbnail_height INTEGER NOT NULL,
  image_thumbnail_width INTEGER NOT NULL,
  image_thumbnail_url TEXT NOT NULL,
  image_source TEXT NOT NULL,
  image_source_url TEXT NOT NULL,
  result_id INTEGER NOT NULL,
  FOREIGN KEY (result_id) REFERENCES results (id) ON DELETE CASCADE
);

-- engine ranks
CREATE TABLE engine_ranks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  engine_name TEXT NOT NULL,
  engine_rank INTEGER NOT NULL,
  engine_page INTEGER NOT NULL,
  engine_on_page_rank INTEGER NOT NULL,
  result_id INTEGER NOT NULL,
  FOREIGN KEY (result_id) REFERENCES results (id) ON DELETE CASCADE
);