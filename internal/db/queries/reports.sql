-- name: CreateDilemmaReport :one
INSERT INTO
    reports (
        reporter_id,
        dilemma_id,
        reason
    )
VALUES ($1, $2, $3)
RETURNING
    *;

-- name: CreateCommentReport :one
INSERT INTO
    reports (
        reporter_id,
        comment_id,
        reason
    )
VALUES ($1, $2, $3)
RETURNING
    *;