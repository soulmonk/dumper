// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: query.sql

package ideas

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countActiveIdeas = `-- name: CountActiveIdeas :one
SELECT COUNT(1) FROM ideas WHERE done_at IS NULL
`

func (q *Queries) CountActiveIdeas(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countActiveIdeas)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createIdea = `-- name: CreateIdea :one
INSERT INTO "ideas" (title, body, created_at)
VALUES ($1, $2, now())
RETURNING id, created_at
`

type CreateIdeaParams struct {
	Title pgtype.Text `json:"title"`
	Body  pgtype.Text `json:"body"`
}

type CreateIdeaRow struct {
	ID        int64            `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) CreateIdea(ctx context.Context, arg CreateIdeaParams) (CreateIdeaRow, error) {
	row := q.db.QueryRow(ctx, createIdea, arg.Title, arg.Body)
	var i CreateIdeaRow
	err := row.Scan(&i.ID, &i.CreatedAt)
	return i, err
}

const doneIdea = `-- name: DoneIdea :one
UPDATE "ideas" SET done_at = now() WHERE id = $1 AND done_at IS NULL
RETURNING done_at
`

func (q *Queries) DoneIdea(ctx context.Context, id int64) (pgtype.Timestamp, error) {
	row := q.db.QueryRow(ctx, doneIdea, id)
	var done_at pgtype.Timestamp
	err := row.Scan(&done_at)
	return done_at, err
}

const getIdea = `-- name: GetIdea :one
SELECT id, title, body, created_at, done_at FROM ideas WHERE id = $1
`

func (q *Queries) GetIdea(ctx context.Context, id int64) (Ideas, error) {
	row := q.db.QueryRow(ctx, getIdea, id)
	var i Ideas
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Body,
		&i.CreatedAt,
		&i.DoneAt,
	)
	return i, err
}

const getIdsOfActiveIdeas = `-- name: GetIdsOfActiveIdeas :many
SELECT id FROM ideas WHERE done_at IS NULL
`

func (q *Queries) GetIdsOfActiveIdeas(ctx context.Context) ([]int64, error) {
	rows, err := q.db.Query(ctx, getIdsOfActiveIdeas)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listIdeas = `-- name: ListIdeas :many
SELECT id, title, body, created_at, done_at FROM ideas ORDER BY created_at DESC
`

func (q *Queries) ListIdeas(ctx context.Context) ([]Ideas, error) {
	rows, err := q.db.Query(ctx, listIdeas)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ideas
	for rows.Next() {
		var i Ideas
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Body,
			&i.CreatedAt,
			&i.DoneAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
