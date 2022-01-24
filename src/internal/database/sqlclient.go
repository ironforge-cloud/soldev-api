package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// GetConnection init a lazy connection with Postgres
func GetConnection() (*sqlx.DB, error) {
	uri := os.Getenv("POSTGRESQL_URL")

	db, err := sqlx.Open("postgres", uri)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

// Close db connection
func Close(db *sql.DB) {
	db.Close()
}
