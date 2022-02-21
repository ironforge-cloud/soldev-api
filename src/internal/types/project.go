package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ProjectCategory struct {
	ID           string         `db:"id"`
	CategoryName string         `db:"category_name"`
	Description  string         `db:"description"`
	Status       string         `db:"status"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	DeletedAt    sql.NullString `db:"deleted_at"`
}

type Project struct {
	ID          string         `db:"id"`
	CategoryID  string         `db:"category_id"`
	ProjectName string         `db:"project_name"`
	Logo        string         `db:"logo"`
	Description string         `db:"description"`
	Status      string         `db:"status"`
	ProjectURL  ProjectURL     `db:"project_url"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	DeletedAt   sql.NullString `db:"deleted_at"`
}

type ProjectsWithCategory struct {
	CategoryName string
	CategoryID   int
	Project      Project
}

type ProjectURL struct {
	Data []URL `db:"data"`
}

type URL struct {
	Name string
	Url  string
}

// Value makes the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a ProjectURL) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan makes the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *ProjectURL) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
