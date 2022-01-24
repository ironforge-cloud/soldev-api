package types

import (
	"database/sql"
	"time"
)

type Bounty struct {
	CompanyID    string         `db:"company_id"`
	ID           string         `db:"id"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	Reward       int            `db:"reward"`
	RewardAsset  string         `db:"reward_asset"`
	Status       string         `db:"status"`
	DeliveryDate string         `db:"delivery_date"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	DeletedAt    sql.NullString `db:"deleted_at"`
}
