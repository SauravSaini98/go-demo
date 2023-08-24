package database

import (
    "database/sql"
    "log"
    "os"
    "fmt"

    _ "github.com/lib/pq" // Import the database driver
)

var db *sql.DB

func InitDatabase() {
    var err error
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	fmt.Println(connStr)

    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Failed to ping the database:", err)
    }
}

func GetDB() *sql.DB {
    return db
}
