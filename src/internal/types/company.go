package types

import (
	"database/sql"
	"time"
)

type Company struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	Logo        string         `db:"logo"`
	BgColor     string         `db:"bg_color"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	DeletedAt   sql.NullString `db:"deleted_at"`
}
