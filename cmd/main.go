package main

import (
	"devstack/internal/db"
	"devstack/internal/handlers"

	"log"
	"net/http"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/users", handlers.GetUsersHandler(database))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
