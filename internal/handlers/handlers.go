package handlers

import (
	"database/sql"
	"devstack/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
