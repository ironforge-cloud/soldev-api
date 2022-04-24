package database

import (
	"api/internal/types"

	"github.com/jmoiron/sqlx"
)

// CreateProject creates a new project in the DB
func CreateProject(db *sqlx.DB, projectCategoryID string, data types.Project) error {

	_, err := db.Exec(`INSERT INTO projects (category_id, project_name, description, logo,
status, project_url) VALUES ($1,$2,$3,$4,$5,$6)`, projectCategoryID, data.ProjectName, data.Description,
		data.Logo, data.Status, data.ProjectURL)

	if err != nil {
		return err
	}

	return nil
}

// UpdateProject update a project in the DB
func UpdateProject(db *sqlx.DB, projectID string, data types.Project) error {

	_, err := db.Exec(`UPDATE projects SET project_name = $2, description = $3, logo = $4,
status = $5, project_url = $6 WHERE  id = $1`, projectID, data.ProjectName, data.Description,
		data.Logo, data.Status, data.ProjectURL)

	if err != nil {
		return err
	}

	return nil
}

// GetProjectByID finds a project by ID
func GetProjectByID(db *sqlx.DB, projectID string) (types.Project, error) {
	var project types.Project

	err := db.Get(&project, `SELECT * FROM projects WHERE id = $1`,
		projectID)

	if err != nil {
		return types.Project{}, err
	}

	return project, nil
}

// GetProjectsGroupedByCategory finds all projects and group them by category.
func GetProjectsGroupedByCategory(db *sqlx.DB) (map[string][]types.ProjectsWithCategory, error) {
	projectsWithCategory := make(map[string][]types.ProjectsWithCategory)

	rows, err := db.Query("SELECT category_name, category_id, project_name, project_url, logo, " +
		"p.description from project_categories FULL JOIN projects p on project_categories.id = p.category_id")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var result types.ProjectsWithCategory
		err = rows.Scan(&result.CategoryName, &result.CategoryID, &result.Project.ProjectName,
			&result.Project.ProjectURL,
			&result.Project.Logo,
			&result.Project.Description)

		projectsWithCategory[result.CategoryName] = append(projectsWithCategory[result.CategoryName],
			result)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return projectsWithCategory, nil
}

// GetProjects finds all projects
func GetProjects(db *sqlx.DB) ([]types.Project, error) {
	var projects []types.Project

	err := db.Select(&projects, `SELECT * from projects WHERE deleted_at is null`)

	if err != nil {
		return nil, err
	}

	return projects, nil
}
