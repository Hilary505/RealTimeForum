package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func Init() {
	dataFolder := "data"
	databaseFile := "forum.db"

	databaseFolderPath := filepath.Join(dataFolder)
	databaseFilePath := filepath.Join(databaseFolderPath, databaseFile)

	if _, err := os.Stat(databaseFolderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(databaseFolderPath, os.ModePerm); err != nil {
			fmt.Println("[DATABASE]: Error creating database folder:", err)
			os.Exit(1)
		}
		fmt.Println("[DATABASE]: Database folder created successfully.")
	}

	if _, err := os.Stat(databaseFilePath); os.IsNotExist(err) {
		dbFile, err := os.Create(databaseFilePath)
		if err != nil {
			fmt.Println("[DATABASE]: Error creating database file:", err)
			os.Exit(1)
		}

		dbFile.Close()
		fmt.Println("[DATABASE]: Database file created successfully.")
	} else {
		fmt.Println("[DATABASE]: Database file already exists. Skipping creation.")
	}

	err := StartDbConnection(databaseFilePath)
	if err != nil {
		fmt.Printf("|starting database connection|%v", err)
		os.Exit(1)
	}
}

var Db *sql.DB

// ==== This function will starting the connection to the database using the SQLite3 driver that works with CGO =====
func StartDbConnection(database_file_path string) error {
	var err error

	Db, err = sql.Open("sqlite3", database_file_path)
	if err != nil {
		return err
	}

	err = Db.Ping()
	if err != nil {
		return err
	}

	if err = CreateUsersTable(Db); err != nil {
		return err
	}

	if err = CreateLikesDislikesTable(Db); err != nil {
		return err
	}

	if err = CreateSessionsTable(Db); err != nil {
		return err
	}

	if err = CreatePostsTable(Db); err != nil {
		return err
	}
	if err = CreateCommentsTable(Db); err != nil {
		return err
	}
    if err = CreatePrivateMessagesTable(Db); err != nil {
		return err
	}
	fmt.Println("[SUCCESS]: Connected to the SQLite database!", nil)
	return nil
}
