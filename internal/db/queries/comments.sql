-- name: CreateComment :one
INSERT INTO
    comments (
        dilemma_id,
        user_id,
        content,
        parent_id
    )
VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: ListCommentsByDilemma :many
SELECT c.id, c.dilemma_id, c.user_id, u.username, c.content, c.parent_id, c.created_at
FROM comments c
    JOIN users u ON u.id = c.user_id
WHERE
    c.dilemma_id = $1
ORDER BY c.created_at DESC
LIMIT $2
OFFSET
    $3;