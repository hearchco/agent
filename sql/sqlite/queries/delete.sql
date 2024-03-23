-- name: DeleteAllResultsWithQuery :exec
DELETE FROM results
WHERE query = ?;

-- name: DeleteAllResultsOlderThanXHours :exec
DELETE FROM results
WHERE created_at < datetime('now', '-1 hour' * ?);

-- name: DeleteAllResultsOlderThanXDays :exec
DELETE FROM results
WHERE created_at < datetime('now', '-1 day' * ?);
