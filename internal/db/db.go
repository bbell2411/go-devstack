package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables(db)

	return db
}

func createTables(db *sql.DB) {
	query := `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	pin INTEGER NOT NULL,
	created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS snippets (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	snippet TEXT NOT NULL,
	language TEXT NOT NULL,
	deleted_at TEXT,
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
