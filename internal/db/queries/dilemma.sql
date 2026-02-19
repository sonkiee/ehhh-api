-- name: CreateDilemma :one
INSERT INTO
    dilemmas (user_id, title, is_anonymous)
VALUES ($1, $2, $3)
RETURNING
    *;

-- name: CreateDilemmaOption :one
INSERT INTO
    dilemma_options (dilemma_id, label)
VALUES ($1, $2)
RETURNING
    *;

-- name: GetDilemma :one
SELECT * FROM dilemmas WHERE id = $1;

-- name: GetDilemmaOptions :many
SELECT * FROM dilemma_options WHERE dilemma_id = $1 ORDER BY id;

-- name: ListFeed :many
SELECT d.id, d.title, d.is_anonymous, d.total_votes, d.created_at, u.username AS author_username
FROM dilemmas d
    JOIN users u ON u.id = d.user_id
WHERE
    d.status = 'active'
ORDER BY d.created_at DESC
LIMIT $1
OFFSET
    $2;