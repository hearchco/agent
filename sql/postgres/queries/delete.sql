-- name: DeleteAllResultsWithQuery :exec
DELETE FROM results
WHERE query = $1;

-- name: DeleteAllResultsOlderThanXHours :exec
DELETE FROM results
WHERE created_at < NOW() - INTERVAL '1 hour' * $1;

-- name: DeleteAllResultsOlderThanXDays :exec
DELETE FROM results
WHERE created_at < NOW() - INTERVAL '1 day' * $1;