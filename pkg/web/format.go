package web

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func formatTimestamp(t pgtype.Timestamp) string {
	if t.Valid {
		return t.Time.UTC().Format(time.RFC3339)
	}
	return ""
}
