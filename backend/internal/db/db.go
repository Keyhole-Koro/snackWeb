package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	// Default to looking for snack.db three levels up from backend_go
	// Adjust as per deployment location
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		// Assuming cwd is .../snackWeb/backend_go
		dbPath = filepath.Join(cwd, "../../../snack.db")
	}

	log.Printf("Connecting to database at: %s", dbPath)

	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
