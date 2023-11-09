-- name: ListIdeas :many
SELECT * FROM ideas ORDER BY created_at DESC;

-- name: CreateIdea :one
INSERT INTO "ideas" (title, body, created_at)
VALUES ($1, $2, now())
RETURNING id, created_at;

-- name: DoneIdea :one
UPDATE "ideas" SET done_at = now() WHERE id = $1 AND done_at IS NULL
RETURNING done_at;

-- name: CountActiveIdeas :one
SELECT COUNT(1) FROM ideas WHERE done_at IS NULL;

-- name: GetIdea :one
SELECT * FROM ideas WHERE id = $1;

-- name: GetIdsOfActiveIdeas :many
SELECT id FROM ideas WHERE done_at IS NULL;
