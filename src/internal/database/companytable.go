package database

import (
	"api/internal/types"
	"time"

	"github.com/jmoiron/sqlx"
)

// CreateCompany saves a company in the DB
func CreateCompany(db *sqlx.DB, data types.Company) error {
	_, err := db.Exec(`INSERT INTO companies(name, description, logo, bg_color) VALUES($1, $2, $3,
$4)`, data.Name, data.Description, data.Logo, data.BgColor)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCompany updates a company in the DB
func UpdateCompany(db *sqlx.DB, companyID string, data types.Company) error {
	_, err := db.Exec(`UPDATE companies SET name = $2, description = $3, logo = $4,
bg_color = $5 WHERE id = $1`, companyID, data.Name, data.Description, data.Logo, data.BgColor)
	if err != nil {

		return err
	}

	return nil
}

// DeleteCompany soft-delete a company from the DB
func DeleteCompany(db *sqlx.DB, companyID string) error {
	_, err := db.Exec(`UPDATE companies SET deleted_at = $2 WHERE id = $1`, companyID, time.Now())
	if err != nil {

		return err
	}

	return nil
}

// GetCompanyByID finds a record in the company's table where id
func GetCompanyByID(db *sqlx.DB, companyID string) (types.Company, error) {
	company := types.Company{}

	err := db.Get(&company, `SELECT * FROM companies WHERE id = $1`,
		companyID)

	if err != nil {
		return types.Company{}, err
	}

	return company, nil
}

// GetAllCompanies finds all companies
func GetAllCompanies(db *sqlx.DB) ([]types.Company, error) {
	var companies []types.Company

	err := db.Select(&companies, `SELECT * FROM companies WHERE deleted_at is null ORDER BY status`)

	if err != nil {
		return nil, err
	}

	return companies, nil
}
