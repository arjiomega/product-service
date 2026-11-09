package product

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatal(err)
	}

	_, err = db.Exec(`TRUNCATE TABLE products CASCADE`)
	if err != nil {
		db.Close()
		t.Fatalf("failed to clean test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}
