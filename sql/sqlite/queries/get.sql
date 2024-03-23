-- name: GetResultsByQueryWithEngineRanks :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsByQueryWithEngineRanksNotOlderThanTimestamp :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND created_at > ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryWithEngineRanks :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryWithEngineRanksNotOlderThanTimestamp :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND created_at > ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsByQueryAndEngineWithEngineRanks :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND engine_ranks.engine_name = ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetResultsByQueryAndEngineWithEngineRanksNotOlderThanTimestamp :many
SELECT * FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND engine_ranks.engine_name = ? AND created_at > ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryAndEngineWithEngineRanks :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND engine_ranks.engine_name = ?
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC;

-- name: GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanTimestamp :many
SELECT * FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = ? AND engine_ranks.engine_name = ? AND created_at > ?
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