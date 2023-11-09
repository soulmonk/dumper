-- name: ListIdeas :many
SELECT * FROM ideas ORDER BY created_at DESC;

-- name: CreateIdea :one
INSERT INTO "ideas" (title, body, created_at)
VALUES ($1, $2, now())
RETURNING id, created_at;

-- name: DoneIdea :one
UPDATE "ideas" SET done_at = now() WHERE id = $1 AND done_at IS NULL
RETURNING done_at;
