-- name: GetResultsByQueryWithEngineRanks :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsByQueryWithEngineRanksNotOlderThanXminutes :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND created_at > datetime('now', '-1 minute' * ?)
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryWithEngineRanks :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryWithEngineRanksNotOlderThanXminutes :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND created_at > datetime('now', '-1 minute' * ?)
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsByQueryAndEngineWithEngineRanks :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND engine_ranks.engine_name = ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND engine_ranks.engine_name = ? AND created_at > datetime('now', '-1 minute' * ?)
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryAndEngineWithEngineRanks :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND engine_ranks.engine_name = ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND engine_ranks.engine_name = ? AND created_at > datetime('now', '-1 minute' * ?)
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsTTLByQuery :one
SELECT created_at FROM results
WHERE query = ?
ORDER BY created_at DESC
LIMIT 1;

-- name: GetResultsTTLByQueryAndEngine :one
SELECT created_at FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND engine_ranks.engine_name = ?
ORDER BY created_at DESC
LIMIT 1;

-- name: AddResult :one
INSERT INTO results (
  query, url, rank, score, title, description
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING id;

-- name: AddImageResult :exec
INSERT INTO image_results (
  image_original_height, image_original_width, image_thumbnail_height, image_thumbnail_width, image_thumbnail_url, image_source, image_source_url, result_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: AddEngineRank :exec
INSERT INTO engine_ranks (
  engine_name, engine_rank, engine_page, engine_on_page_rank, result_id
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: DeleteAllResultsWithQuery :exec
DELETE FROM results
WHERE query = ?;

-- name: DeleteAllResultsOlderThanXHours :exec
DELETE FROM results
WHERE created_at < datetime('now', '-1 hour' * ?);

-- name: DeleteAllResultsOlderThanXDays :exec
DELETE FROM results
WHERE created_at < datetime('now', '-1 day' * ?);
