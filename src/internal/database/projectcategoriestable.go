package database

import (
	"api/internal/types"

	"github.com/jmoiron/sqlx"
)

// CreateProjectCategory creates a new project category in the DB
func CreateProjectCategory(db *sqlx.DB, data types.ProjectCategory) error {
	_, err := db.Exec(`INSERT INTO project_categories(category_name, description,
status) VALUES($1, $2, $3)`, data.CategoryName, data.Description, data.Status)
	if err != nil {
		return err
	}

	return nil
}
