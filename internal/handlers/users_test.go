package handlers

import (
	"database/sql"
	"devstack/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
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

func TestGetUsersHandler_NoRows(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	_, err := db.Exec(`DELETE FROM users`)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
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

func TestGetUsersHandler_MethodNotAllowed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodDelete, "/users/", nil)

	w := httptest.NewRecorder()

	handler := GetUsersHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)

}

func TestGetUserHandler(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)

	w := httptest.NewRecorder()
	handler := GetUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	var user models.Users
	err := json.NewDecoder(res.Body).Decode(&user)
	require.NoError(t, err)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, 1, user.ID)
}

func TestGetUserHandler_IDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/users/99787", nil)
	w := httptest.NewRecorder()

	handler := GetUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetUserHandler_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/users/hello", nil)
	w := httptest.NewRecorder()

	handler := GetUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetUserHandler_MethodNotAllowed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()

	handler := GetUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)

}

func TestCreateUserHandler(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	body := `{
		"Name":"Bell",
		"Pin":"1234"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users/", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler := CreateUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var user models.Users
	json.NewDecoder(res.Body).Decode(&user)

	assert.Equal(t, "Bell", user.Name)
	assert.Equal(t, 1234, user.Pin)
	assert.NotZero(t, user.ID)

}

func TestCreateUserHandler_InvalidName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	body := `{
		"Name":"   ",
		"Pin":"1234"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users/", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler := CreateUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateUserHandler_InvalidPin(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	body := `{
		"Name":"bell",
		"Pin":"123564"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users/", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler := CreateUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateUserHandler_InvalidPinEmpty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	body := `{
		"Name":"bell",
		"Pin":""
	}`

	req := httptest.NewRequest(http.MethodPost, "/users/", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler := CreateUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateUserHandler_InvalidMethod(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	body := `{
		"Name":"bell",
		"Pin":"1234"
	}`

	req := httptest.NewRequest(http.MethodGet, "/users/", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler := CreateUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestCreateUserHandler_NoBody(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodPost, "/users/", nil)
	w := httptest.NewRecorder()
	handler := CreateUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateUserHandler_MissingFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	body := `{
		"Name":"bell",
	}`

	req := httptest.NewRequest(http.MethodPost, "/users/", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler := CreateUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateUserHandler_InvalidFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	body := `{
		"Name":"bell",
		"Pin":"1234",
		"extra":"973"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users/", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler := CreateUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateUserHandler_InvalidPinData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	body := `{
		"Name":"bell",
		"Pin":"abc"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users/", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler := CreateUserHandler(db)
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
