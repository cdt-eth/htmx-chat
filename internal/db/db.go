package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() error {
    // Connect to Postgres
    // Create tables
    // Set up migrations
} 