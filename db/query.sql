-- name: ListIdeas :many
SELECT * FROM ideas;

-- name: CreateIdea :one
INSERT INTO "ideas" (title, body, created_at)
VALUES ($1, $2, now())
RETURNING id, created_at;
