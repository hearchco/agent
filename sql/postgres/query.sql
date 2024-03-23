-- name: GetResultsByQueryWithEngineRanks :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsByQueryWithEngineRanksNotOlderThanXminutes :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND created_at > NOW() - INTERVAL '1 minute' * $2
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryWithEngineRanks :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryWithEngineRanksNotOlderThanXminutes :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND created_at > NOW() - INTERVAL '1 minute' * $2
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsByQueryAndEngineWithEngineRanks :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND engine_ranks.engine_name = $2
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND engine_ranks.engine_name = $2 AND created_at > NOW() - INTERVAL '1 minute' * $3
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryAndEngineWithEngineRanks :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND engine_ranks.engine_name = $2
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND engine_ranks.engine_name = $2 AND created_at > NOW() - INTERVAL '1 minute' * $3
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsTTLByQuery :one
SELECT created_at FROM results
WHERE query = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetResultsTTLByQueryAndEngine :one
SELECT created_at FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND engine_ranks.engine_name = $2
ORDER BY created_at DESC
LIMIT 1;

-- name: AddResult :one
INSERT INTO results (
  query, url, rank, score, title, description
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id;

-- name: AddImageResult :exec
INSERT INTO image_results (
  image_original_height, image_original_width, image_thumbnail_height, image_thumbnail_width, image_thumbnail_url, image_source, image_source_url, result_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: AddEngineRank :exec
INSERT INTO engine_ranks (
  engine_name, engine_rank, engine_page, engine_on_page_rank, result_id
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: DeleteAllResultsWithQuery :exec
DELETE FROM results
WHERE query = $1;

-- name: DeleteAllResultsOlderThanXHours :exec
DELETE FROM results
WHERE created_at < NOW() - INTERVAL '1 hour' * $1;

-- name: DeleteAllResultsOlderThanXDays :exec
DELETE FROM results
WHERE created_at < NOW() - INTERVAL '1 day' * $1;