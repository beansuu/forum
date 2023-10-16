package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Initialize() error {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}

	createPostsTable := `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err = db.Exec(createPostsTable)
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
