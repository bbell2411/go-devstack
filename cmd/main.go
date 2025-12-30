package main

import (
	"devstack/internal/db"

	"log"
	"net/http"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	mux := http.NewServeMux()

	// handlers.RegisterRoutes(mux, database)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
