package main

import (
	"log"
	"net/http"
	"database/sql"
	"my_project/database"
	"my_project/api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func loadEnvironments() {
    envErr := godotenv.Load()
    if envErr != nil {
        log.Fatal("Error loading .env file")
    }
}

func init() {
    loadEnvironments()
    database.InitDatabase()
}

func main() {
    // Database connection parameters
    http.HandleFunc("/api/v1/users", api.HandleUsersApi)
	http.ListenAndServe(":8080", nil)
}
