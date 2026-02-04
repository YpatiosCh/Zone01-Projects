package database

import (
	"database/sql"
	"fmt"
	uuid "forum/pkg/UUID"
	"os"
)

var (
	schemaPath = GetEnv("FORUM_SCHEMA_PATH", "./internal/database/schema.sql")
	dbPath     = GetEnv("FORUM_DB_PATH", "./forum.db")
)

func GetEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// createSchema creates the database schema
func createSchema(db *sql.DB) error {
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	// Convert bytes to string
	schema := string(schemaBytes)

	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}

	fmt.Println("Schema created successfully")
	return nil
}

// populateCategories populates the categories table
func populateCategories(db *sql.DB) error {
	categories := []string{
		"General Discussion",
		"Suggestions",
		"Announcements",
		"Off-Topic",
		"Tech News",
	}

	for _, category := range categories {
		// create uuid for each category
		id := uuid.GenerateUUID()
		_, err := db.Exec("INSERT INTO category (id, name) VALUES (?, ?)", id, category)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %v", category, err)
		}
	}

	fmt.Println("Categories populated successfully")
	return nil
}

// InitDB initializes the database only if it does not exist
func InitDB() (*sql.DB, error) {

	// Check if the database file exists
	dbExists := true
	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			dbExists = false
		} else {
			return nil, fmt.Errorf("failed to check database file: %v", err)
		}
	} else if fileInfo.Size() == 0 {
		// File exists but is empty
		dbExists = false
	}

	// Connect to SQLite database
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)

	// If the database didn't exist before, create schema and populate data
	if !dbExists {
		fmt.Println("Initializing new database...")

		if err := createSchema(db); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to create tables: %v", err)
		}

		if err := populateCategories(db); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to populate categories: %v", err)
		}
		fmt.Println("Database initialized successfully")
	} else {
		fmt.Println("Database already exists. Skipping initialization.")
	}

	return db, nil
}
