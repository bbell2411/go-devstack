package handlers

import (
	"database/sql"
	"devstack/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		INSERT INTO users (name, pin) VALUES ("Test User2", 1234);

	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestGetUsersHandler(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler := GetUsersHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	var users []models.Users
	err := json.NewDecoder(res.Body).Decode(&users)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(users))

	assert.Equal(t, "Test User", users[0].Name)
	assert.Equal(t, "Test User2", users[1].Name)
}

func TestGetUsersHandlerNOROWS(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	_, err := db.Exec(`DELETE FROM users`)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler := GetUsersHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	var users []models.Users
	err = json.NewDecoder(res.Body).Decode(&users)
	assert.NoError(t, err)

	assert.Len(t, users, 0)
}

func TestGetUsersHandlerINVALIDMETHOD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodDelete, "/users", nil)

	w := httptest.NewRecorder()

	handler := GetUsersHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)

}
