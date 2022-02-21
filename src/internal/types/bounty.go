package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Bounty struct {
	CompanyID    string         `db:"company_id"`
	ID           string         `db:"id"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	Reward       int            `db:"reward"`
	RewardAsset  string         `db:"reward_asset"`
	Tags         Tags           `db:"tags"`
	Status       string         `db:"status"`
	DeliveryDate string         `db:"delivery_date"`
	URL          string         `db:"url"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	DeletedAt    sql.NullString `db:"deleted_at"`
}

type Tags struct {
	Names []string `json:"names"`
}

type BountyStats struct {
	TotalBalance  int `db:"total_balance"`
	PaidBalance   int `db:"paid_balance"`
	TotalBounties int `db:"total_bounties"`
}

// Value makes the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a Tags) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan makes the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *Tags) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
