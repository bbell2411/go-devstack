package handlers

import (
	"database/sql"
	"devstack/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		rows, err := db.Query(`
		SELECT id, name, created_at
		FROM users`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		fmt.Println("in the handler")
		users := []models.Users{}

		for rows.Next() {
			var user models.Users

			if err := rows.Scan(
				&user.ID,
				&user.Name,
				&user.CreatedAt,
			); err != nil {
				log.Printf("SCAN ERROR: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/users/")
		if idStr == "" {
			http.Error(w, "Missing ID", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var user models.Users

		err = db.QueryRow(`
		select id, name, created_at
		from users
		where id= ?
		`, id).Scan(
			&user.ID,
			&user.Name,
			&user.CreatedAt,
		)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}

}
