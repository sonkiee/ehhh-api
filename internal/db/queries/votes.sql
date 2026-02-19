-- name: CreateVoteTx :one

WITH
    inserted AS (
        INSERT INTO
            votes (
                user_id,
                dilemma_id,
                option_id
            )
        VALUES ($1, $2, $3)
        RETURNING
            option_id
    ),
    upd_opt AS (
        UPDATE dilemma_options o
        SET
            vote_count = vote_count + 1
        WHERE
            o.id = (
                SELECT option_id
                FROM inserted
            )
        RETURNING
            o.dilemma_id
    ),
    upd_d AS (
        UPDATE dilemmas d
        SET
            total_votes = total_votes + 1
        WHERE
            d.id = (
                SELECT dilemma_id
                FROM upd_opt
            )
        RETURNING
            d.id,
            d.total_votes
    )
SELECT (
        SELECT id
        FROM upd_d
    ) AS dilemma_id, (
        SELECT total_votes
        FROM upd_d AS total_votes
    );