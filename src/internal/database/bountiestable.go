package database

import (
	"api/internal/types"
	"time"

	"github.com/jmoiron/sqlx"
)

// CreateBounty creates a new bounty in the DB
func CreateBounty(db *sqlx.DB, companyID string, data types.Bounty) error {

	_, err := db.Exec(`INSERT INTO bounties (company_id, title, description, reward, reward_asset,
status, delivery_date) VALUES ($1,$2,$3,$4,$5,$6,$7)`, companyID, data.Title, data.Description,
		data.Reward, data.RewardAsset, data.Status, data.DeliveryDate)

	if err != nil {
		return err
	}

	return nil
}

// UpdateBounty update a bounty in the DB
func UpdateBounty(db *sqlx.DB, bountyID string, data types.Bounty) error {

	_, err := db.Exec(`UPDATE bounties SET title = $2, description = $3, reward = $4,
reward_asset = $5, status = $6, delivery_date = $7 WHERE  id = $1`, bountyID,
		data.Title, data.Description, data.Reward, data.RewardAsset, data.Status, data.DeliveryDate)

	if err != nil {
		return err
	}

	return nil
}

// DeleteBounty soft-delete a company from the DB
func DeleteBounty(db *sqlx.DB, bountyID string) error {
	_, err := db.Exec("UPDATE bounties SET deleted_at = $2 WHERE id = $1", bountyID, time.Now())
	if err != nil {

		return err
	}

	return nil
}

// GetAllBountiesByCompanyID finds all bounties for a project
func GetAllBountiesByCompanyID(db *sqlx.DB, projectID string) ([]types.Bounty, error) {
	var bounties []types.Bounty

	err := db.Select(&bounties, `SELECT * from bounties WHERE company_id = $1 AND deleted_at is null
`, projectID)

	if err != nil {
		return nil, err
	}

	return bounties, nil
}

// GetBountyByID finds a bounty by ID
func GetBountyByID(db *sqlx.DB, bountyID string) (types.Bounty, error) {
	var bounty types.Bounty

	err := db.Get(&bounty, `SELECT * FROM bounties WHERE id = $1`,
		bountyID)

	if err != nil {
		return types.Bounty{}, nil
	}

	return bounty, nil
}
