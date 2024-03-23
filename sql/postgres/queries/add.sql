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