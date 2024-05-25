package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// database connection
func DatabaseConn() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db_host := os.Getenv("DB_HOST")
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	connString := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", db_user, db_name, db_pass, db_host, db_port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Error while connecting to db", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		log.Fatal("Error while pinging db", err)
	}

	return db

}

func CloseDbConn(db *sql.DB) {
	db.Close()
}
