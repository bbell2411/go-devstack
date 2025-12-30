package handlers

import (
	"database/sql"
	// "net/http"
	// "net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			pin INTEGER NOT NULL,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP
		);

		INSERT INTO users (name, pin) VALUES ("Test User", 1234);
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}
