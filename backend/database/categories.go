package database

import (
	"database/sql"
	"fmt"
)

func CreateCategoriesTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("nil database connection")
	}

	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	// Insert default categories if they don't exist
	defaultCategories := []string{
		"General",
		"Technology",
		"Sports",
		"Entertainment",
		"Science",
		"Health",
		"Education",
		"Business",
	}

	for _, category := range defaultCategories {
		_, err := db.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", category)
		if err != nil {
			return err
		}
	}

	return nil
} 